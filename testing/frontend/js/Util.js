//Util
const getElemenetById = (id) => {
  return document.getElementById(id);
};

const getQuerySelector = (selector) => {
  return document.querySelector(selector);
};

const makeQR = (value) => {
  return new QRCode("qrcode", {
    text: value,
    width: 250,
    height: 250,
    colorDark: "#000000",
    colorLight: "#ffffff",
    correctLevel: QRCode.CorrectLevel.H,
  });
};

const eraseAvailableQR = () => {
  let qrCodeContainer = getElemenetById("qrcode");
  qrCodeContainer.title = "";
  while (qrCodeContainer.hasChildNodes()) {
    qrCodeContainer.removeChild(qrCodeContainer.firstChild);
  }
};

const getCallingType = (type) => {
  switch (type.toLowerCase()) {
    case "call":
      return "call";
    default:
      console.log("Calling type is not a right format " + String(type));
      return null;
  }
};

// const processSMAddress = (address) => {
//   if (typeof address !== "string") {
//     console.log("Address is not a right format " + String(address));
//     return null;
//   }
//   return address.slice(2);
// };

const formatAddressToHex = (address) => {
  if (address === null) return;
  let prefix = "0";
  // return prefix.repeat(64 - address.length) + address.toLowerCase();
  // console.log(prefix.repeat(64 - address.length) + address)

  return prefix.repeat(64 - address.length) + address;
};

const processAmount = (amount) => {
  amountHex = Number(amount).toString(16);
  console.log('amount:',amountHex)
  return amountHex;
};

const formatNumberToHex = (number) => {
  numberHex = Number(number).toString(16);
  let prefix = "0";
  return prefix.repeat(64 - numberHex.length) + String(numberHex);
};

const processInput = (input) => {
  let formattedInput = "";
  input.forEach((element) => {
    switch (element.type) {
      case "hash":
        formattedInput += element.value;
        break;
      case "num":
        formattedInput += formatNumberToHex(element.value);
        break;
      case "address":
        formattedInput += formatAddressToHex(element.value);
       
        break;
      default:
        return;
    }
  });
  return formattedInput;
};

const formatInput = (_type, _contractAddress, _amount, _input) => {
  let type = getCallingType(_type);
  let contractAddress = _contractAddress;
  let amount = processAmount(_amount);
  let input = processInput(_input);

  console.log(`${type}|${contractAddress}|${amount}|${input}`);
  return `${type}|${contractAddress}|${amount}|${input}`;
};
