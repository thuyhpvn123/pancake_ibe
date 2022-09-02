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
  addLiquidityMessage.type = "AddLiquidity";
  addLiquidityMessage.message = walletAddress;

  sendMessage(addLiquidityMessage);

  //Format input
  let inputMessage = structuredClone(SMRouter.addLiquidity);
  inputMessage.parameter[1].value = SMAToken.address;
  inputMessage.parameter[2].value = SMBToken.address;
  inputMessage.parameter[3].value = getElemenetById(`token-A-input`).value;
  inputMessage.parameter[4].value = getElemenetById(`token-B-input`).value;
  inputMessage.parameter[7].value = wallAddress;
  inputMessage.parameter[8].value = DEADLINE;

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
