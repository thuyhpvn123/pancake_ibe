//Client
const handleApproveToken = (tokenName) => {
    let amountToken = getElemenetById(`token-${tokenName}-input`).value;
    if (amountToken.length === 0) return;
  
    let smartContract;
    switch (tokenName) {
      case "a":
      case "A":
        tokenName = "A";
        smartContract = SMAToken; //make a clone for smartContract
        break;
      case "b":
      case "B":
        tokenName = "B";
        smartContract = SMBToken;
        break;
      default:
        console.log("Invalid token");
        return;
    }
  
    if (tokenName === "") return;
    console.log(`Call approve Token ${tokenName}`);
  
    let inputMessage = structuredClone(smartContract.approve);
    //Assign value to SMContract Parameter
    inputMessage.parameter[1].value = SMRouter.address; //spender
    inputMessage.parameter[2].value = amountToken; //amount value;
  
    eraseAvailableQR();
    makeQR(
      formatInput("call", smartContract.address, "", inputMessage.parameter)
    );
  };
  
  getElemenetById("approve-a-btn").addEventListener("click", () => {
    handleApproveToken("A");
  });
  
  getElemenetById("approve-b-btn").addEventListener("click", () => {
    handleApproveToken("B");
  });
  
  const handleSupply = () => {
    //Send to backend a flag that this address is calling liquidity adding
    addLiquidityMessage = structuredClone(messageForm);
    addLiquidityMessage.type = "swap";
    addLiquidityMessage.message = walletAddress;
  
    sendMessage(addLiquidityMessage);
  
    //Format input
    // let inputMessage = structuredClone(SMRouter.addLiquidity);
    // inputMessage.parameter[1].value = SMAToken.address;
    // inputMessage.parameter[2].value = SMBToken.address;
    // inputMessage.parameter[3].value = getElemenetById(`token-A-input`).value;
    // inputMessage.parameter[4].value = getElemenetById(`token-B-input`).value;
    // inputMessage.parameter[7].value = wallAddress;
    // inputMessage.parameter[8].value = DEADLINE;

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
    //Print QRCode
    eraseAvailableQR();
    makeQR(formatInput("call", SMRouter.address, "", inputMessage.parameter));
  };
  
  getElemenetById("supply-btn").addEventListener("click", handleSupply);
  
  const handleResetQR = () => {
    eraseAvailableQR();
  };
  
  getElemenetById("reset-qr-btn").addEventListener("click", handleResetQR);
  
  getElemenetById("simple-contract-btn").addEventListener("click", () => {
    let inputMessage = structuredClone(SMSimpleContract.numberChange);
    inputMessage.parameter[1].value = "2";
    eraseAvailableQR();
  
    makeQR(
      formatInput("call", SMSimpleContract.address, "", inputMessage.parameter)
    );
  });
  