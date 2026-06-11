// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract PRAEscrow is Ownable, ReentrancyGuard {
    IERC20 public praToken;
    
    enum TradeStatus { Pending, Paid, Completed, Cancelled }
    
    struct Trade {
        uint256 id;
        address buyer;
        address seller;
        uint256 praAmount;
        uint256 irrAmount;      // مبلغ ریالی توافق شده (برای شفافیت)
        TradeStatus status;
        uint256 createdAt;
        uint256 paidAt;
        uint256 completedAt;
    }
    
    uint256 public tradeCounter;
    mapping(uint256 => Trade) public trades;
    mapping(address => uint256[]) public userTrades;
    
    event TradeCreated(uint256 indexed tradeId, address indexed buyer, address indexed seller, uint256 praAmount, uint256 irrAmount);
    event PaymentConfirmed(uint256 indexed tradeId, address indexed buyer);
    event TradeCompleted(uint256 indexed tradeId, address indexed seller);
    event TradeCancelled(uint256 indexed tradeId);
    
    constructor(address _praToken) {
        praToken = IERC20(_praToken);
    }
    
    // ایجاد معامله جدید (فروشنده توکن را قفل می‌کند)
    function createTrade(uint256 _praAmount, uint256 _irrAmount) external nonReentrant {
        require(_praAmount > 0, "PRA amount must be > 0");
        require(_irrAmount > 0, "IRR amount must be > 0");
        
        // انتقال توکن از فروشنده به قرارداد
        require(praToken.transferFrom(msg.sender, address(this), _praAmount), "Token transfer failed");
        
        tradeCounter++;
        trades[tradeCounter] = Trade({
            id: tradeCounter,
            buyer: address(0),
            seller: msg.sender,
            praAmount: _praAmount,
            irrAmount: _irrAmount,
            status: TradeStatus.Pending,
            createdAt: block.timestamp,
            paidAt: 0,
            completedAt: 0
        });
        
        userTrades[msg.sender].push(tradeCounter);
        
        emit TradeCreated(tradeCounter, address(0), msg.sender, _praAmount, _irrAmount);
    }
    
    // خریدار معامله را قبول می‌کند و متعهد به پرداخت IRR می‌شود (آفلاین)
    function acceptTrade(uint256 _tradeId) external nonReentrant {
        Trade storage trade = trades[_tradeId];
        require(trade.status == TradeStatus.Pending, "Trade not pending");
        require(trade.buyer == address(0), "Trade already has buyer");
        require(trade.seller != msg.sender, "Seller cannot be buyer");
        
        trade.buyer = msg.sender;
        trade.status = TradeStatus.Paid; // فرض می‌کنیم پرداخت IRR خارج از زنجیره انجام شده
        
        userTrades[msg.sender].push(_tradeId);
        
        emit PaymentConfirmed(_tradeId, msg.sender);
    }
    
    // تأیید نهایی توسط ادمین / اربیتر (پس از تأیید دریافت IRR)
    function confirmPaymentAndRelease(uint256 _tradeId) external onlyOwner nonReentrant {
        Trade storage trade = trades[_tradeId];
        require(trade.status == TradeStatus.Paid, "Trade not in paid state");
        require(trade.buyer != address(0), "No buyer assigned");
        
        // انتقال توکن به خریدار
        require(praToken.transfer(trade.buyer, trade.praAmount), "Transfer to buyer failed");
        
        trade.status = TradeStatus.Completed;
        trade.completedAt = block.timestamp;
        
        emit TradeCompleted(_tradeId, trade.seller);
    }
    
    // لغو معامله (فقط در حالت Pending و توسط فروشنده)
    function cancelTrade(uint256 _tradeId) external nonReentrant {
        Trade storage trade = trades[_tradeId];
        require(trade.status == TradeStatus.Pending, "Cannot cancel");
        require(msg.sender == trade.seller, "Only seller can cancel");
        
        // برگرداندن توکن به فروشنده
        require(praToken.transfer(trade.seller, trade.praAmount), "Refund failed");
        
        trade.status = TradeStatus.Cancelled;
        
        emit TradeCancelled(_tradeId);
    }
    
    // دریافت اطلاعات معامله
    function getTrade(uint256 _tradeId) external view returns (
        uint256 id,
        address buyer,
        address seller,
        uint256 praAmount,
        uint256 irrAmount,
        TradeStatus status,
        uint256 createdAt
    ) {
        Trade storage trade = trades[_tradeId];
        return (
            trade.id,
            trade.buyer,
            trade.seller,
            trade.praAmount,
            trade.irrAmount,
            trade.status,
            trade.createdAt
        );
    }
    
    // دریافت لیست معاملات یک کاربر
    function getUserTrades(address _user) external view returns (uint256[] memory) {
        return userTrades[_user];
    }
}
