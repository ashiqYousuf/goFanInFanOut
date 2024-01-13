## Go Fanning in and Fanning Out Pattern

 
 Lets suppose you have a random number generator Go-Routine & you have to filter out Primes from those random numbers & then you want to log say K primes only, So We have in this Case following functions


 ### Generator (Write)------>(Read)   PrimeFinder   (Write)------>(Read)   LogKPrimes 

 ### Where ----------> Represents a communication channel with corresponding reader & writer

 In order to sync all three go-routines we have to make sure to share some communication pipe (channels) with all those 3 routines as shown above. Now This program might take sometime as this is highly inefficient (Yeah!)

 Since we know primes have no relation with each other [This is key point here], if P1 is a prime Number, and P2 is also a prime, they have no relation with each other, means P1 being Prime doesn't impact the Primality of other Numbers.

 So why not spin multiple instances of the PrimeFinder to quickly filter Primes & then Log them.

                         
### Generator ---->FanOut   [PrimeFinder1 PrimeFinder2 PrimeFinder3]  FanIn----------> LogKPrimes
                         

