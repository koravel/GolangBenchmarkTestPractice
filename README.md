# GolangBenchmarkTestPractice
A little practice on Golang. These tests represent matrices adding with 3 different scenarios of L1 CPU cache usage:
1. MatrixCombination
   Plain matrices addition with 1:1 indexes scenario
   ![image](https://github.com/koravel/GolangBenchmarkTestPractice/assets/26851016/7f272668-dc58-4701-b054-4061df831f99)
3. MatrixReversedCombination
   Transposed matrix addition
   ![image](https://github.com/koravel/GolangBenchmarkTestPractice/assets/26851016/5b5e2d94-662d-4a53-9eb3-39c087e98ee2)
4. MatrixReversedCombinationPerBlock
   Transposed matrix addition with breaking down processing to n batches
   ![image](https://github.com/koravel/GolangBenchmarkTestPractice/assets/26851016/e360b65d-5220-4d7f-b7d8-5f00f02cf3f9)
Based on this article: ![Go and CPU caches](https://teivah.medium.com/go-and-cpu-caches-af5d32cc5592)

Was tested on Ryzen 5700X(16 logical cores) and 64GB RAM
As a result most efficient was 64 items batch', giving 50+% boost just from awaring how CPU caching works.
Might be useful in some edge/hi-load cases.
