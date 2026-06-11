import { ethers } from 'ethers';
import BasketABI from './BasketABI.json';

const contract = new ethers.Contract(contractAddress, BasketABI, signer);

// خواندن اطلاعات
const allocations = await contract.getCurrentAllocations();

// نوشتن تراکنش
const tx = await contract.deposit(amountIn, minAmountOut);
await tx.wait();