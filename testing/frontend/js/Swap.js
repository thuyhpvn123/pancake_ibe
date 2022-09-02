let priceAToB;

priceAToB = 2;

let priceBToA = 1 / priceAToB;

handleOnChange = (tokenName) => {
  switch (tokenName.toUpperCase()) {
    case "A":
      calculateTokenAToB();
      console.log("Calculate A");
      break;
    case "B":
      calculateTokenBToA();
      break;
    default:
      alert("Input wrong token name");
      break;
  }
};

let calculateTokenAToB = () => {
  let tokenAAmount = getElemenetById("token-A-input").value;
  let tokenBAmount = tokenAAmount * priceAToB;
  getElemenetById("token-B-input").value = tokenBAmount;
};

let calculateTokenBToA = () => {
  let tokenBAmount = getElemenetById("token-A-input").value;
  let tokenAAmount = tokenBAmount * priceBToA;
  getElemenetById("token-A-input").value = tokenAAmount;
};

let showTokenCurrency = () => {
  getElemenetById("price-content").value = priceAToB;
};

let getPriceAOnB = () => {
  console.log("Get Price A On B");
  getPriceMessage = {
    type: "GetPriceList",
    message: "",
  };

  let tokenAAmount = getElemenetById("token-A-input").value;
  if (tokenAAmount.length == 0) {
    tokenAAmount = 1;
  }
  getPriceMessage.message = tokenAAmount;
  sendMessage(getPriceMessage);
};

getElemenetById("price-btn").addEventListener("click", getPriceAOnB);

socket.onmessage = (msg) => {
  console.log(JSON.parse(msg.data));
  data = JSON.parse(msg.data);
  output.innerHTML += "Server: " + msg.data + "\n";
  console.log(data.type);
  switch (data.type) {
    case "GetPriceList":
      getElemenetById("token-B-input").value = data.message;
      break;
    default:
      break;
  }
};

const handleSupply = () => {
  //Send to backend a flag that this address is calling liquidity adding
  addLiquidityMessage = structuredClone(messageForm);
  addLiquidityMessage.type = "swap";
  addLiquidityMessage.message = walletAddress;

  sendMessage(addLiquidityMessage);


  let inputMessage = structuredClone(SMRouter.swap);
  inputMessage.parameter[1].value = getElemenetById(`token-A-input`).value;
  inputMessage.parameter[4].value = wallAddress;
  inputMessage.parameter[5].value = DEADLINE;
  inputMessage.parameter[7].value = SMAToken.address;
  inputMessage.parameter[8].value = SMBToken.address;

//   38ed1739
// 00000000000000000000000000000000000000000000000000000002540be400  10000000000
// 0000000000000000000000000000000000000000000000000000000005f5e100  100000000
// 00000000000000000000000000000000000000000000000000000000000000a0 160 line 5
// 000000000000000000000000fee8665978caf2e902a24b4b100613883ffc4d2f to address
// 0000000000000000000000000000000000000000000000056bc75e2d63100000 time 100000000000000000000
// 0000000000000000000000000000000000000000000000000000000000000002 do dai mang
// 000000000000000000000000E6fBE813230f087813c35c950FC46e3bee4847D1 token1
// 000000000000000000000000f1b5dc17F84FC6e0fA632BF81406748ABfb6F6Cd token2
// call|3FFb75AcB68A021e978b8bEc4E4762d32060152f||
// 38ed1739
// 0000000000000000000000000000000000000000000000000000000000000001
// 0000000000000000000000000000000000000000000000000000000000000000
// 00000000000000000000000000000000000000000000000000000000000000a0
// 000000000000000000000000355308801ba607d9224370aefe9bfbbe54989991
// 000000000000000000000000000000000000000000000001007a33ee6d770000
// 0000000000000000000000000000000000000000000000000000000000000002
// 000000000000000000000000b9C40c5054333975e4fEE5b2972f2481422CD48D
// 000000000000000000000000148b3c21920A743625974Bf7E7f3b8D675534b74
// 38ed1739
// 0000000000000000000000000000000000000000000000000000000000000000
// 0000000000000000000000000000000000000000000000000000000000000000
// 00000000000000000000000000000000000000000000000000000000000000a0
// 000000000000000000000000d85ae9a6ef6185aea70b1b18c3d3bfd1253ea74e
// 000000000000000000000000000000000000000000000001007a33ee6d770000
// 0000000000000000000000000000000000000000000000000000000000000002
// 000000000000000000000000b9C40c5054333975e4fEE5b2972f2481422CD48D
// 000000000000000000000000148b3c21920A743625974Bf7E7f3b8D675534b74
  //Print QRCode
  eraseAvailableQR();
  makeQR(formatInput("call", SMRouter.address, "", inputMessage.parameter));
};

getElemenetById("supply-btn").addEventListener("click", handleSupply);
let sendMessage = (msg) => {
  console.log(msg);
  socket.send(JSON.stringify(msg));
};

