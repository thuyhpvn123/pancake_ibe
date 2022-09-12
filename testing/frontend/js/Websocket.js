var flag =1;

const output = document.getElementById("log-content");
let approveToken="";
let balanceToken="";
const socket = new WebSocket("ws://127.0.0.1:3000/ws");
const walletAddress = getElemenetById("wallet-id").innerHTML;
// const walletAddressOwner = getElemenetById("wallet-id-owner").innerHTML;

let socketActive = false;

console.log("Imported");
// * Websocket
// Connect to server successfully
const messageForm = {
  type: "",
  message: "",
};

socket.onopen = (msg) => {
  socketActive = true;

  output.innerHTML += "Status: Connected\n";

  //Send walletAddress to server
  let walletMessage = structuredClone(messageForm);
  walletMessage.type = "WalletMessage";
  walletMessage.message = walletAddress;

  sendMessage(walletMessage);
};

// WS connection's closed
socket.onclose = (event) => {
  console.log("WS Connection is closed: ", event);
};

// WS connection having errors
socket.onerror = (error) => {
  console.log("Socket Error: ", error);
};
socket.onmessage = (msg) => {
// 9  let data1 = msg.data;
  const data12 = JSON.parse(msg.data);
  // console.log("data12:",data12.message)
  output.innerHTML += "Server: " + msg.data + "\n";
  // console.log(data12.type);
  switch (data12.type) {
    case "GetPriceList":
      if(flag==1){
        console.log('111');
        getElemenetById("token-B-input").value = data12.message;

      } else{
        console.log('222');
        getElemenetById("token-A-input").value = data12.message;
        flag=1;
      }
        break;
    case "GetApprove":
      switch (approveToken){
        case "A":
          if(parseInt(getElemenetById(`token-A-input`).value )<= parseInt(data12.message)){
            document.getElementById('appstatus-a').innerHTML = 'Approved'

          }else{
            document.getElementById('appstatus-a').innerHTML = 'Not yet'
          }
          break;
        case "B":
          if(parseInt(getElemenetById(`token-B-input`).value )<= parseInt(data12.message)){
            document.getElementById('appstatus-b').innerHTML = '<span class="fa fa-plus mr-2"></span>Approved'
            console.log(getElemenetById(`token-B-input`).value)
            console.log(data12.message)

          }else{
            document.getElementById('appstatus-b').innerHTML = '<span class="fa fa-plus mr-2"></span>Not yet'
          }
          break;
        case "LP":
          if(parseInt(getElemenetById('liquidity-input').value) <= parseInt(data12.message)){
            document.getElementById('appstatus-lp').innerHTML = 'Approved'
    
          }else{
            document.getElementById('appstatus-lp').innerHTML = 'Not yet'
          }
          break;
        default:
          break;
      }
      break;
      case "GetBalance":
        switch (balanceToken){
          case "A":
              document.getElementById('balanceA').innerHTML = data12.message
            break;
          case "B":
              document.getElementById('balanceB').innerHTML = data12.message
            break;
          case "LP":
            document.getElementById('balanceLP').innerHTML = data12.message
          break;

          default:
            break;
        }
        break;
  
    default:
      break;
  }
};


const wallAddress = getElemenetById("wallet-id").innerHTML;


//Choose Token Address
const $tokenA = document.getElementById('tokenA');
const $createResultA = document.getElementById('create-resultA');
const $tokenB = document.getElementById('tokenB');
const $createResultB = document.getElementById('create-resultB');
const $LPtoken = document.getElementById('LPtoken');
const $createResultLP = document.getElementById('create-resultLP');
const $tokenMTDA = document.getElementById('tokenMTDA');
const $tokenMTDB = document.getElementById('tokenMTDB');
let tokenAddressA ="";
let tokenAddressB ="";
let lpToken ="";
$tokenA.addEventListener('submit', async(e) => {
  e.preventDefault();
  document.getElementById('appstatus-a').innerHTML =""
  document.getElementById('appstatus-lp').innerHTML =""
  var name,flag =1
  name = $('#tokenAName').val()

  if( name ==''){
    flag=0
    $('.error_name').html("Please type token address")
  }else{
    $('.error_name').html("")
  }
  if(flag==1){
    try{
      console.log("getBalance111a")
      tokenAddressA = name;
      await getBalanceToken(name); 
      console.log("getBalance222a")   
      $createResultA.innerHTML = name;
      balanceToken = "A";
    }catch(e){
      console.log(e)
    $createResultA.innerHTML = `Ooops... there was an error while trying to get token address`;
    }
  }
$tokenA.reset()

});
$tokenB.addEventListener('submit', async(e) => {
  e.preventDefault();
  document.getElementById('appstatus-b').innerHTML =""
  document.getElementById('appstatus-lp').innerHTML =""
  var name,flag =1
  name = $('#tokenBName').val()

  if( name ==''){
    flag=0
    $('.error_name').html("Please type token address")
  }else{
    $('.error_name').html("")
  }
  if(flag==1){
    try{
      tokenAddressB = name;
      console.log("getBalance111b")
      await getBalanceToken(name);
      console.log("getBalance222b")
      $createResultB.innerHTML = name;
      balanceToken ="B";
    }catch(e){
      console.log(e)
    $createResultB.innerHTML = `Ooops... there was an error while trying to get token address`;
    }
  }
$tokenB.reset()

});

$tokenMTDB.addEventListener("click", async(e) => {
  e.preventDefault();
  document.getElementById('appstatus-b').innerHTML =""
  document.getElementById('balanceB').innerHTML =""
  document.getElementById('appstatus-lp').innerHTML =""
  tokenAddressB = MTDToken.address;
  // getBalanceToken(tokenAddressB);
  $createResultB.innerHTML = MTDToken.address;
  balanceToken = "LP";
  
});
$tokenMTDA.addEventListener("click", async(e) => {
  console.log("thuy day")
  e.preventDefault();
  document.getElementById('appstatus-a').innerHTML =""
  document.getElementById('balanceA').innerHTML =""
  document.getElementById('appstatus-lp').innerHTML =""
  console.log(MTDToken)
  tokenAddressA = MTDToken.address;
  // getBalanceToken(tokenAddressA);
  $createResultA.innerHTML = MTDToken.address;
  balanceToken = "LP";
});
$LPtoken.addEventListener('submit', async(e) => {
  e.preventDefault();
  document.getElementById('appstatus-lp').innerHTML =""
  var name,flag =1
  name = $('#tokenLPName').val()

  if( name ==''){
    flag=0
    $('.error_name').html("Please type token address")
  }else{
    $('.error_name').html("")
  }
  if(flag==1){
    try{
      lpToken= name;
      await getBalanceToken(name);
      $createResultLP.innerHTML = name;
      balanceToken ="LP";
    }catch(e){
      console.log(e)
    $createResultB.innerHTML = `Ooops... there was an error while trying to get token address`;
    }
  }
$LPtoken.reset()

});

//GetPriceList   

const getBalanceToken = (address) => {
  console.log("getBalance333")
    let getBalanceMessage = {
      type: "GetBalance",
      message: `${address}`,  
    }
    // console.log(getPriceMessage);
    sendMessage(getBalanceMessage);
    console.log("getBalance444")
};

// let getPrice = (event) => {
//   if (event.key == "Enter") {
//     event.preventDefault();
//     let getPriceMessage = {
//       type: "GetPriceList",
//       message: "",
//     };

//     switch (event.target.id) {
//       case "token-A-input":
//         // console.log("GetPriceList from token MTD for DDT");
//         getPriceMessage.message = getElemenetById("token-A-input").value;
//         // approveToken = "A";
//         break;
//       case "token-B-input":
//         flag=0;
//         // console.log("GetPriceList from token DDT for MTD");
//         getPriceMessage.message = getElemenetById("token-B-input").value;
//         // approveToken = "B";
//         break;
//       default:
//         break;
//     }

//     // console.log(getPriceMessage);
//     sendMessage(getPriceMessage);
//   }
// };
let getPrice = (event) => {
  if (event.key == "Enter") {
    event.preventDefault();
    let getPriceMessage = {
      type: "GetPriceList",
      message: "",
    };

    switch (event.target.id) {
      case "token-A-input":
        getPriceMessage.message= getElemenetById("token-A-input").value+','+tokenAddressA+','+tokenAddressB;
        break;
      case "token-B-input":
        flag=0;
        getPriceMessage.message= getElemenetById("token-B-input").value+','+tokenAddressB+','+tokenAddressA;
        break;
      default:
        break;
    }
    sendMessage(getPriceMessage);
  }
};

getElemenetById("token-A-input").addEventListener("keypress", (event)=>{
  if ( tokenAddressA=== tokenAddressB){
    alert(' Can not swap same token address. Choose another token address!')
  }else{
    getPrice(event)
  }
});
getElemenetById("token-B-input").addEventListener("keypress", (event)=>{
  if (tokenAddressA=== tokenAddressB){
    alert(' Can not swap same token address. Choose another token address!')
  }else{
    getPrice(event)
  }
});

//Approve Button
// const handleApproveBtn = () => {
//   eraseAvailableQR();
//   let inputMessage;
//   switch (approveToken.toUpperCase()) {
//     case "A":
//       inputMessage = structuredClone(SMAToken.approve);
//       inputMessage.parameter[1].value = SMRouter.address;
//       inputMessage.parameter[2].value = getElemenetById("token-A-input").value;
//       console.log(inputMessage);
//       makeQR(formatInput("call", tokenAddressA, "", inputMessage.parameter));
//       break;
//     case "B":
//       inputMessage = structuredClone(SMBToken.approve);
//       inputMessage.parameter[1].value = SMRouter.address;
//       inputMessage.parameter[2].value = getElemenetById("token-B-input").value;
//       makeQR(formatInput("call", tokenAddressB, "", inputMessage.parameter));
//       break;
//     default:
//       console.log("Approve token error");
//       break;
//   }
// };

// getElemenetById("approve-btn").addEventListener("click", handleApproveBtn);
//
// let priceAToB;

// priceAToB = 2;

// let priceBToA = 1 / priceAToB;

// handleOnChange = (tokenName) => {
//   switch (tokenName.toUpperCase()) {
//     case "A":
//       calculateTokenAToB();
//       console.log("Calculate A");
//       break;
//     case "B":
//       calculateTokenBToA();
//       break;
//     default:
//       alert("Input wrong token name");
//       break;
//   }
// };

// let calculateTokenAToB = () => {
//   let tokenAAmount = getElemenetById("token-A-input").value;
//   let tokenBAmount = tokenAAmount * priceAToB;
//   getElemenetById("token-B-input").value = tokenBAmount;
// };

// let calculateTokenBToA = () => {
//   let tokenBAmount = getElemenetById("token-A-input").value;
//   let tokenAAmount = tokenBAmount * priceBToA;
//   getElemenetById("token-A-input").value = tokenAAmount;
// };

// let showTokenCurrency = () => {
//   getElemenetById("price-content").value = priceAToB;
// };


const handleAdd = () => {
  //Send to backend a flag that this address is calling liquidity adding
  addMessage = structuredClone(messageForm);
  addMessage.type = "addliquidity";
  addMessage.message = walletAddress;

  sendMessage(addMessage);

  //Format input
  let inputMessage 
  let amount 
  if(tokenAddressA == MTDToken.address){
    inputMessage = structuredClone(SMRouter.addLiquidityMTD);
    inputMessage.parameter[1].value = tokenAddressB;
    inputMessage.parameter[2].value = getElemenetById(`token-B-input`).value;
    inputMessage.parameter[5].value = wallAddress;
    inputMessage.parameter[6].value = DEADLINE;
    amount = getElemenetById(`token-A-input`).value
    console.log("amount:",amount)
  }else if(tokenAddressB == MTDToken.address){
    inputMessage = structuredClone(SMRouter.addLiquidityMTD);
    inputMessage.parameter[1].value = tokenAddressA;
    inputMessage.parameter[2].value = getElemenetById(`token-A-input`).value;
    inputMessage.parameter[5].value = wallAddress;
    inputMessage.parameter[6].value = DEADLINE;
    amount = getElemenetById(`token-B-input`).value
    console.log("amount:",amount)
  }else{
    inputMessage= structuredClone(SMRouter.addLiquidity);
    inputMessage.parameter[1].value = tokenAddressA;
    inputMessage.parameter[2].value = tokenAddressB;
    inputMessage.parameter[3].value = getElemenetById(`token-A-input`).value;
    inputMessage.parameter[4].value = getElemenetById(`token-B-input`).value;
    inputMessage.parameter[7].value = wallAddress;
    inputMessage.parameter[8].value = DEADLINE;
  
  }

  //Print QRCode
  eraseAvailableQR();
  makeQR(formatInput("call", SMRouter.address, amount, inputMessage.parameter));
};
const handleSwap = () => {
  //Send to backend a flag that this address is calling liquidity adding
  swapMessage = structuredClone(messageForm);
  swapMessage.type = "swap";
  swapMessage.message = walletAddress;

  sendMessage(swapMessage);
  let inputMessage;
  let amount;
  if(tokenAddressA == MTDToken.address)
  {
    inputMessage = structuredClone(SMRouter.swapMTD);
    inputMessage.parameter[3].value = wallAddress;
    inputMessage.parameter[4].value = DEADLINE;
    inputMessage.parameter[6].value = MTDToken.address;
    inputMessage.parameter[7].value = tokenAddressB;
    amount = getElemenetById(`token-A-input`).value
  }
  else if(tokenAddressB == MTDToken.address)
  {
    inputMessage = structuredClone(SMRouter.swapToMTD);
    inputMessage.parameter[1].value = getElemenetById(`token-A-input`).value;
    inputMessage.parameter[4].value = wallAddress;
    inputMessage.parameter[5].value = DEADLINE;
    inputMessage.parameter[7].value = tokenAddressA;
    inputMessage.parameter[8].value = MTDToken.address;
    amount = getElemenetById(`token-B-input`).value
  }
  else{
    inputMessage = structuredClone(SMRouter.swap);
    inputMessage.parameter[1].value = getElemenetById(`token-A-input`).value;
    inputMessage.parameter[4].value = wallAddress;
    inputMessage.parameter[5].value = DEADLINE;
    inputMessage.parameter[7].value = tokenAddressA;
    inputMessage.parameter[8].value = tokenAddressB;
  }

  //Print QRCode
  eraseAvailableQR();
  makeQR(formatInput("call", SMRouter.address, amount, inputMessage.parameter));
};
const handleRemove = () => {
  //Send to backend a flag that this address is calling liquidity adding
  removeMessage = structuredClone(messageForm);
  removeMessage.type = "remove";
  removeMessage.message = walletAddress;

  sendMessage(removeMessage);

  //Format input
  let inputMessage = structuredClone(SMRouter.removeLiquidity);
  inputMessage.parameter[1].value = tokenAddressA;
  inputMessage.parameter[2].value = tokenAddressB;
  inputMessage.parameter[3].value = getElemenetById(`liquidity-input`).value;
  inputMessage.parameter[6].value = wallAddress;
  inputMessage.parameter[7].value = DEADLINE;

  //Print QRCode
  eraseAvailableQR();
  makeQR(formatInput("call", SMRouter.address, "", inputMessage.parameter));
};

const handleApproveTokenA = () => {
  //Format input
  let inputMessage = structuredClone(SMToken.approve);
  //Assign value to SMContract Parameter
  inputMessage.parameter[1].value = SMRouter.address; //spender
  inputMessage.parameter[2].value = getElemenetById(`token-A-input`).value; //amount value;

  //Print QRCode
  eraseAvailableQR();
  makeQR(formatInput("call", tokenAddressA, "", inputMessage.parameter));
};
const handleApproveTokenB = () => {
  //Format input
  let inputMessage = structuredClone(SMToken.approve);
  //Assign value to SMContract Parameter
  inputMessage.parameter[1].value = SMRouter.address; //spender
  inputMessage.parameter[2].value = getElemenetById(`token-B-input`).value; //amount value;

  //Print QRCode
  eraseAvailableQR();
  makeQR(formatInput("call", tokenAddressB, "", inputMessage.parameter));
};
const handleApproveLPToken = () => {
  //Format input
  let inputMessage = structuredClone(LPToken.approve);
  //Assign value to SMContract Parameter
  inputMessage.parameter[1].value = SMRouter.address; //spender
  inputMessage.parameter[2].value = getElemenetById(`liquidity-input`).value; //amount value;

  //Print QRCode
  eraseAvailableQR();
  makeQR(formatInput("call", LPToken.address, "", inputMessage.parameter));
};

// let  handleAppStatus= (event) => {
//   event.preventDefault();
//   let getApproveMessage = {
//     type: "GetApprove",
//     message: "",
//   };

//   switch (event.target.id) {
//     case "appstatus-a-btn":
//       getApproveMessage.message = tokenAddressA;
//       approveToken = "A";
//       break;
//     case "appstatus-b-btn":
//       getApproveMessage.message = tokenAddressB;
//       approveToken = "B";
//       break;
//     case "appstatus-lp-btn":
      
//       getApproveMessage.message = lpToken;
//       approveToken = "LP";
//       break;
  
//     default:
//       break;
//   }

//   // console.log(getPriceMessage);
//   sendMessage(getApproveMessage);

// };


getElemenetById("supply-btn").addEventListener("click",handleSwap);
getElemenetById("add-btn").addEventListener("click",handleAdd);
getElemenetById("remove-btn").addEventListener("click",handleRemove);
getElemenetById("approve-a-btn").addEventListener("click",handleApproveTokenA);
getElemenetById("approve-b-btn").addEventListener("click",handleApproveTokenB);
getElemenetById("approve-b-btn").addEventListener("click",handleApproveTokenB);
getElemenetById("approve-lp-btn").addEventListener("click",handleApproveLPToken);
getElemenetById("appstatus-a-btn").addEventListener("click",handleAppStatus);
getElemenetById("appstatus-b-btn").addEventListener("click",handleAppStatus);
getElemenetById("appstatus-lp-btn").addEventListener("click",handleAppStatus);
let sendMessage = (msg) => {
  console.log(msg);
  socket.send(JSON.stringify(msg));
};
getElemenetById("reset-qr-btn").addEventListener("click", handleResetQR);
const handleResetQR = () => {
  eraseAvailableQR();
};
//addLiquidityMTD value=10000
// f305d719
// 000000000000000000000000b9C40c5054333975e4fEE5b2972f2481422CD48D token
// 0000000000000000000000000000000000000000000000000000000000002710 amountdesired
// 0000000000000000000000000000000000000000000000000000000000000001 amountTokenMin
// 0000000000000000000000000000000000000000000000000000000000000001 amountMTDMin
// 000000000000000000000000d85ae9a6ef6185aea70b1b18c3d3bfd1253ea74e to
// 0000000000000000000000000000000000000000000000056bc75e2d63100000 time

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

// //swapExactETHForToken 
// 7ff36ab5
// 0000000000000000000000000000000000000000000000000000000000000001  amountOutMin
// 0000000000000000000000000000000000000000000000000000000000000080  128 line4
// 0000000000000000000000001fa4ad1d255980ff7e7578b36382ef0488889131  to address
// 0000000000000000000000000000000000000000000000056bc75e2d63100000  time 100000000000000000000
// 0000000000000000000000000000000000000000000000000000000000000002  do dai mang
// 0000000000000000000000003589d24a8b038c85118a43e81065d72ca009d949  wbnb
// 000000000000000000000000636af98334aa7f53bcd6d1bb67fed1a802a7e180  token2
// // swapExactTokensForETH amountIn=100000
// // 
// 18cbafe5
// 00000000000000000000000000000000000000000000000000000000000186a0  100000
// 0000000000000000000000000000000000000000000000000000000000000001  AmountOutMin
// 00000000000000000000000000000000000000000000000000000000000000a0  160 line5
// 0000000000000000000000001fa4ad1d255980ff7e7578b36382ef0488889131  to address
// 0000000000000000000000000000000000000000000000056bc75e2d63100000  time 100000000000000000000
// 0000000000000000000000000000000000000000000000000000000000000002  do dai mang
// 000000000000000000000000636af98334aa7f53bcd6d1bb67fed1a802a7e180  token1
// 0000000000000000000000003589d24a8b038c85118a43e81065d72ca009d949  wbnb
// remove 10000000000000000000
// baa2abde
// 000000000000000000000000E6fBE813230f087813c35c950FC46e3bee4847D1 token1
// 000000000000000000000000f1b5dc17F84FC6e0fA632BF81406748ABfb6F6Cd token2
// 0000000000000000000000000000000000000000000000008ac7230489e80000 liquidity
// 0000000000000000000000000000000000000000000000008ac7230489e80000 amountAmin
// 0000000000000000000000000000000000000000000000008ac7230489e80000 amountBmin
// 000000000000000000000000bdda90332da8c4ea6f27aea75a8b12d14b770293 to
// 0000000000000000000000000000000000000000000000056bc75e2d63100000 deadtime
// approve 
// 095ea7b3
// 000000000000000000000000aA557dafC14C7b84E37d479036D1630773FCc788
// 000000000000000000000000000000000001ed09bead87c0378d8e6400000000
// balanceOf 
// 70a08231
// 0000000000000000000000004e131e2a811f967977bc35f3159b590163e06dfb
