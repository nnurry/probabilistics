# PROBABILISTICS

This repository is a compilation of flexible probabilistic data structures and algorithms, based on popular implementations. It leaves users with choices to replace hash functions, and extend its functionalities. 

Another version of [my old PDSA repository](https://github.com/nnurry/pds-old).

Note: 
- This is a research repo, not intended for production
- /v1 is comprised of sketch implementations of:
    + Hash functions (64-bit Murmur3)
    + Hash schemers (Kir-Mit/EDH) 
    + 2^k-counter
    + Classic/naive-counting/counting bloom filters
    + Basic/stochastic average probabilistic counting
    + LogLog/SuperLogLog/HyperLogLog (all 64-bits as opposed to 32-bit in past researches since I only implemented 64-bit hash functions)
    