## Consideration

The calculation for quality is done line by line, rather than at the end of reading all data points.

The complexity is O(n) for checking the quality of thermometers. 
The complexity can be adjusted to be O(1) in the best case for checking the quality for humidity, as once a data point is out of the 1% range, the whole dataset should be regarded as `discard`. And, in the worst case is o(n) to be `ok`, it needs to read all records. 

The Data points along with their timestamp can be stored in the dictionary as well but have been removed, as it is unnecessary to capture that base on the given requirements. 

## run 
The log path is `sample.log`, `reader.js` read the logs and print on the screen. The recommended Node Version is 17.

```
node reader.js
```

## Description

37Widgets makes inexpensive thermometers and humidity sensors. In order to spot check the manufacturing process, some units are put in a test environment (for an unspecified amount of time) and their readings are logged. The test environment has a static, known temperature and relative humidity, but the sensors are expected to fluctuate a bit.                                                                                                          
 
As a developer, your task is to process the logs and automate the quality control evaluation. The evaluation criteria are as follows:
 
1) For a thermometer, it is branded “ultra precise” if the mean of the readings is within 0.5 degrees of the known temperature, and the standard deviation is less than 3. It is branded “very precise” if the mean is within 0.5 degrees of the room, and the standard deviation is under 5. Otherwise, it’s sold as “precise`.                  
 
2) For a humidity sensor, it must be discarded unless it is within 1% of the reference value for all readings.  
 
An example log looks like the following. The first line means that the room was held at a constant 70 degrees, 45% relative humidity. Subsequent lines either identify a sensor (<type> <name>) or give a reading (<time> <name> <value>.                                                
 
```
reference 70.0 45.0                                    
thermometer temp-1                                      
2007-04-05T22:00 temp-1 72.4                            
2007-04-05T22:01 temp-1 76.0                            
2007-04-05T22:02 temp-1 79.1                            
2007-04-05T22:03 temp-1 75.6                            
2007-04-05T22:04 temp-1 71.2                            
2007-04-05T22:05 temp-1 71.4                            
2007-04-05T22:06 temp-1 69.2                            
2007-04-05T22:07 temp-1 65.2                            
2007-04-05T22:08 temp-1 62.8                            
2007-04-05T22:09 temp-1 61.4                            
2007-04-05T22:10 temp-1 64.0                            
2007-04-05T22:11 temp-1 67.5                            
2007-04-05T22:12 temp-1 69.4                            
thermometer temp-2                                      
2007-04-05T22:01 temp-2 69.5                            
2007-04-05T22:02 temp-2 70.1                            
2007-04-05T22:03 temp-2 71.3                            
2007-04-05T22:04 temp-2 71.5                            
2007-04-05T22:05 temp-2 69.8                            
humidity hum-1                                          
2007-04-05T22:04 hum-1 45.2                            
2007-04-05T22:05 hum-1 45.3                            
2007-04-05T22:06 hum-1 45.1                            
humidity hum-2                                          
2007-04-05T22:04 hum-2 44.4                            
2007-04-05T22:05 hum-2 43.9                            
2007-04-05T22:06 hum-2 44.9                            
2007-04-05T22:07 hum-2 43.8                            
2007-04-05T22:08 hum-2 42.1 
``` 

``` 
Output                                                  
temp-1: precise                                        
temp-2: ultra precise                                  
hum-1: OK                                              
hum-2: discard   
```                                       
 
The log should be read from stdin. In the end, you will own this process, so you should solve the problem as described, but feel free to advocate for any changes you think would make sense to improve the process (split into multiple files, change log format, etc).
