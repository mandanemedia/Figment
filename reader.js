const fs = require('fs');
const readline = require('readline');

// input file
const logPath = './sample.log';

// look up Keywords on the logs
const referenceKeyword = 'reference';
const deviceType1 = 'thermometer';
const deviceType2 = 'humidity';
const refernces = [];
// dictionary
const devices = {};

const detectReferences = (line, referenceKeyword) => {
    // remove reference keyword from first line
    const firstLine = line.replace(referenceKeyword);

    // split them base on whitespace and grab index 1 and 2 
    const array = firstLine.split(' ');
    if(array.length < 3){
        throw `Please set reference value in proper format`
    }

    refernces.push(Number(array[1])) 
    refernces.push(Number(array[2]))
    refernces.forEach(current=>{  
        if(typeof current === NaN){
            throw "Reference value is not set properly"
        }
    });
    return refernces;
}

const addDeviceName = (line) => {
    if(line.indexOf(deviceType1) !== -1){
        // split them base on whitespace and grab index 1 and 2 
        const array = line.split(' ');
        if(array.length < 2){
            throw `Please set Device Name in proper format for ${deviceType1}`;
        }
        // there is a duplication device name, can throw error here
        if( !devices[array[1]] ){
            devices[array[1]] = { 
                type: deviceType1,
                data:[],
                quality:'', 
                sum : 0, 
                n:0,
                sumSquareRoot:0,
                mean:0,
                sd:0,
                // if the mean of the readings is within 0.5 degrees of the known temperature
                checkMeanWithin: true 
            };
        }
        else {
            throw `Duplication device name detected`;
        }
        return array[1];
    }
    else{
        // split them base on whitespace and grab index 1 and 2 
        const array = line.split(' ');
        if(array.length < 2){
            throw `Please set Device Name in proper format for ${deviceType2}`;
        }
        // there is a duplication device name, can throw error here
        if( !devices[array[1]] ){
            devices[array[1]] = { 
                type: deviceType2,
                data:[],
                quality:'ok', // default value 
                discarded: false 
            };
        }
        else {
            throw `Duplication device name detected`;
        }
        return array[1];
    }
}

const checkDataPoint = (line, lastDevice) => {
    const data = line.split(' '+lastDevice+' ');
    if(data.length < 2){
        throw `Please set data in proper format for ${deviceType2}`;
    }
    if(typeof Number(data[1]) === NaN){
        throw `value for ${lastDevice} is not a number at ${data[0]}`;
    }
    else{
        // there is no need to capture the whole data set 
        // this can be enable, but commented out to reduce space complexity
        // devices[lastDevice].data.push({time:data[0], value:Number(data[1])});
        // devices[lastDevice].data.push(Number(data[1]));

        if(devices[lastDevice].type === deviceType1){
            checkThermometerQuality(devices[lastDevice], Number(data[1]));
        }else{
            checkHumidityQuality(devices[lastDevice], Number(data[1]));
        }
    }
}

const checkThermometerQuality = (device, currentValue) => {
    device.sum += currentValue;
    device.n += 1;
    device.mean = (device.sum/device.n);
    // sum((current - avg)^2)
    device.sumSquareRoot += Math.pow(currentValue-device.mean, 2);
    device.sd = Math.sqrt(device.sumSquareRoot/(device.n-1))

    // if the mean of the readings is within 0.5 degrees of the refernce    
    if(  device.mean-0.5 <= refernces[0]  && refernces[0] <= device.mean+0.5) {
        device.checkMeanWithin = true;
    }
    if( device.sd < 3 && device.checkMeanWithin){
        device.quality = 'ultra precise'
    }
    else if( device.sd < 5 && device.checkMeanWithin){
        device.quality = 'very precise'
    }
    else {
        device.quality = 'precise'
    }
}

const checkHumidityQuality = (device, currentValue) => {
    // if discarded is already true ignore checking it
    if(!device.discarded){
        const upper = refernces[1]*1.01;
        const lower = refernces[1]*0.99;
        if(!(lower <= currentValue && currentValue <= upper)) {
            device.discarded = true;
            device.quality = 'discard'     
        }
    }
}

const processLineByLine = async () => {
  let lastDevice = '';
  try {
    const fileStream = fs.createReadStream(logPath);
    const rl = readline.createInterface({
        input: fileStream,
        crlfDelay: Infinity
    });
    
    for await (const line of rl) {
        // detect the first line  
        if( refernces.length === 0 ){
                if(line.indexOf(referenceKeyword)=== -1){
                    throw "Log file is not in proper format on first line"
                }
                detectReferences(line, referenceKeyword);
        }

        // detect devices' name
        else if(line.indexOf(deviceType1) !== -1 || line.indexOf(deviceType2) !== -1 ){
            lastDevice = addDeviceName(line);
        }

        // detect if last device name is in line, it is device data 
        else if(line.indexOf(lastDevice)){
            checkDataPoint(line, lastDevice);
        }
        else{
            throw "Mismatch in the log file"
        }

    }
  } catch (err) {
    console.error('Rejected by', err);
  } finally{
    // console.log(`devices:${JSON.stringify(devices)}`);
    // print the result
    console.log(`Output`);
    for (const [key, value] of Object.entries(devices)) {
        console.log(`${key}: ${value.quality}`);
    }
  }

}

processLineByLine();
