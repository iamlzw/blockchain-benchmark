# blockchain-benchmark
a hyperledger fabric benchmark tool

### how to use 

### 1、clone git repository
```bash
$ git clone 
```

### configure 
you need replace ``crypto-config` directory which include your crypto material
#### confgure e2e_config.yaml
this file was use to connect fabric network

#### configure benchmark.yaml 
this file was use to set goroutine nums and test case nums,

### run and generate benchmark result
run main.go 

after about 1 minute ,you will get invoke.html and query.html
![image](https://user-images.githubusercontent.com/27334218/135966708-d523e862-f8ec-4146-b7fc-97c7415cf621.png)

![image](https://user-images.githubusercontent.com/27334218/135966754-1ad0092b-8a22-426b-8331-ccc6623dfe79.png)


