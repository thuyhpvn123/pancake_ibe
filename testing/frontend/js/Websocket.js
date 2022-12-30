
var flag =1;
var output = document.getElementById("log-content");
var approveToken="";
var balanceToken="";
var socket = new WebSocket("ws://127.0.0.1:3000/ws");
var walletAddress = getElemenetById("wallet-id").innerHTML;
// var walletAddressOwner = getElemenetById("wallet-id-owner").innerHTML;

var socketActive = false;

console.log("Imported");
// * Websocket
// Connect to server successfully
var messageForm = {
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
  fetch("/compoundToken");
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
  var data12 = JSON.parse(msg.data);
  output.innerHTML += "Server: " + msg.data + "\n";
  switch (data12.type) {
    case "GetPriceList":
      if(flag==1){
        getElemenetById("token-B-input").value = data12.message;

      } else{
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
              document.getElementById('balanceA').innerHTML = data12.message[0]
            break;
          case "B":
              document.getElementById('balanceB').innerHTML = data12.message[0]
            break;
          case "LP":
            document.getElementById('balanceLP').innerHTML = data12.message[0]
          break;
          case "MTDA":
              document.getElementById('balanceA').innerHTML = data12.message
            break;
          case "MTDB":
            document.getElementById('balanceB').innerHTML = data12.message
          break;

          default:
            break;
        }
        break;
      case "GetFarmPoolInfo":
        loadDataCompound();
      //   fetch("/staking");
        // loadData(data12.message);
  
    default:
      break;
  }
  // console.log("message:",data12.message.event)
  // console.log("event day ne:", data12.message.event)
  // console.log("event day ne:", typeof data12.message.event)
  switch (data12.message.event) {
    case "Deposit":
      console.log("hiiiiiii", data12.message.data.pid)
      var a = document.getElementById(`liquidity_${data12.message.data.pid}`).value
     console.log("aaaaa", a);
     var b =BigInt(a)+BigInt(data12.message.data.amount);
     console.log("bbbb", b);
     console.log("data12.message.data.amount", data12.message.data.amount)
     document.getElementById(`liquidity_${data12.message.data.pid}`).value =b;
     document.getElementById(`liquidity_${data12.message.data.pid}`).innerHTML =b
     updateAPR(data12.message.data.pid);
     break;
    case "Withdraw":
      console.log("hiaaaaaaa")
      var c = document.getElementById(`liquidity_${data12.message.data.pid}`).value
      var d =BigInt(c)-BigInt(data12.message.data.amount);
      document.getElementById(`liquidity_${data12.message.data.pid}`).value =d;
      document.getElementById(`liquidity_${data12.message.data.pid}`).innerHTML =d
      updateAPR(data12.message.data.pid);
      break;
     
  }
};


var wallAddress = getElemenetById("wallet-id").innerHTML;


//Choose Token Address
var $tokenA = document.getElementById('tokenA');
var $createResultA = document.getElementById('create-resultA');
var $tokenB = document.getElementById('tokenB');
var $createResultB = document.getElementById('create-resultB');
var $LPtoken = document.getElementById('LPtoken');
var $createResultLP = document.getElementById('create-resultLP');
var $tokenMTDA = document.getElementById('tokenMTDA');
var $tokenMTDB = document.getElementById('tokenMTDB');
var tokenAddressA ="";
var tokenAddressB ="";
var lpToken ="";
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
      $createResultA.innerHTML = name;
      tokenAddressA = name;
      await getBalanceToken(name);       
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
      $createResultB.innerHTML = name;
      tokenAddressB = name;
      console.log("getBalance111b")
      await getBalanceToken(name);
      console.log("getBalance222b")
      
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
  getBalanceToken(tokenAddressB);
  $createResultB.innerHTML = MTDToken.address;
  balanceToken = "MTDB";
  
});

$tokenMTDA.addEventListener("click", async(e) => {
  // console.log("thuy day")
  e.preventDefault();
  document.getElementById('appstatus-a').innerHTML =""
  document.getElementById('balanceA').innerHTML =""
  document.getElementById('appstatus-lp').innerHTML =""
  console.log(MTDToken)
  tokenAddressA = MTDToken.address;
  getBalanceToken(tokenAddressA);
  $createResultA.innerHTML = MTDToken.address;
  balanceToken = "MTDA";
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

  
//get balance of token
var getBalanceToken = (address) => {
  console.log("getBalance333")
  var getBalanceMessage = {
      type: "GetBalance",
      message: `${address},${wallAddress}`,  
    }
    // console.log(getPriceMessage);
    sendMessage(getBalanceMessage);
    console.log("getBalance444")
};
var flagComp =0
var getCompoundTab = () => {

  if (flagComp == 0){
    console.log("getCompoundTab111")
    var getCompoundTabMessage = {
        type: "GetCompoundTab",
        message: `${wallAddress}`,  
      }
      // console.log(getPriceMessage);
      sendMessage(getCompoundTabMessage);
      console.log("getCompoundTab444")
  
  }
  flagComp =1
};
var getStakingTab = () => {
  console.log("getStakingTab111")
  let getFarmTable = structuredClone(messageForm);
  getFarmTable.type = "GetFarmPoolInfo";
  getFarmTable.message = "GetFarmPoolInfo";
  sendMessage(getFarmTable);
  console.log("getStakingTab444")
  
};


//GetPriceList 
var getPrice = (event) => {
  if (event.key == "Enter") {
    event.preventDefault();
    var getPriceMessage = {
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


var handleAdd = () => {
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

var handleSwap = () => {
  //Send to backend a flag that this address is calling liquidity adding
  swapMessage = structuredClone(messageForm);
  swapMessage.type = "swap";
  swapMessage.message = walletAddress;

  sendMessage(swapMessage);
  var inputMessage;
  var amount;
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
var handleRemove = () => {
  //Send to backend a flag that this address is calling liquidity adding
  removeMessage = structuredClone(messageForm);
  removeMessage.type = "remove";
  removeMessage.message = walletAddress;

  sendMessage(removeMessage);

  //Format input
  var inputMessage = structuredClone(SMRouter.removeLiquidity);
  inputMessage.parameter[1].value = tokenAddressA;
  inputMessage.parameter[2].value = tokenAddressB;
  inputMessage.parameter[3].value = getElemenetById(`liquidity-input`).value;
  inputMessage.parameter[6].value = wallAddress;
  inputMessage.parameter[7].value = DEADLINE;

  //Print QRCode
  eraseAvailableQR();
  makeQR(formatInput("call", SMRouter.address, "", inputMessage.parameter));
};

var handleApproveTokenA = () => {
  //Format input
  var inputMessage = structuredClone(SMToken.approve);
  //Assign value to SMContract Parameter
  inputMessage.parameter[1].value = SMRouter.address; //spender
  inputMessage.parameter[2].value = getElemenetById(`token-A-input`).value; //amount value;

  //Print QRCode
  eraseAvailableQR();
  makeQR(formatInput("call", tokenAddressA, "", inputMessage.parameter));
};
var handleApproveTokenB = () => {
  //Format input
  var inputMessage = structuredClone(SMToken.approve);
  //Assign value to SMContract Parameter
  inputMessage.parameter[1].value = SMRouter.address; //spender
  inputMessage.parameter[2].value = getElemenetById(`token-B-input`).value; //amount value;

  //Print QRCode
  eraseAvailableQR();
  makeQR(formatInput("call", tokenAddressB, "", inputMessage.parameter));
};
var handleApproveLPToken = () => {
  //Format input
  var inputMessage = structuredClone(LPToken.approve);
  //Assign value to SMContract Parameter
  inputMessage.parameter[1].value = SMRouter.address; //spender
  inputMessage.parameter[2].value = getElemenetById(`liquidity-input`).value; //amount value;

  //Print QRCode
  eraseAvailableQR();
  makeQR(formatInput("call", LPToken.address, "", inputMessage.parameter));
};

var  handleAppStatus= (event) => {
  event.preventDefault();
  var getApproveMessage = {
    type: "GetApprove",
    message: "",
  };

  switch (event.target.id) {
    case "appstatus-a-btn":
      getApproveMessage.message = `${tokenAddressA},${wallAddress}`;
      approveToken = "A";
      break;
    case "appstatus-b-btn":
      getApproveMessage.message = `${tokenAddressB},${wallAddress}`;
      approveToken = "B";
      break;
    case "appstatus-lp-btn":
      
      getApproveMessage.message = `${lpToken},${wallAddress}`;
      approveToken = "LP";
      break;
  
    default:
      break;
  }

  // console.log(getPriceMessage);
  sendMessage(getApproveMessage);

};
var cakePerblock = 1* (10**18);

var updateAPR = async(idPool) => {
  console.log("updateAPR")
  var APR ,lpReward,totalValueOfCake;
    var farmBaseReward,totalFee,volume24h;
  var priceOfCake = 20;//getAmountOut in router????
  var cakePerYear=cakePerblock*60*60*24*365;
  liquidity =document.getElementById(`liquidity_${idPool}`).value
  multiplier = document.getElementById(`multiplier_${idPool}`).value
  console.log("liquidity moi ne:",liquidity)
  volume24h = 10000000000000000// dang fix cung theo usd
  totalValueOfCake= cakePerYear *priceOfCake;
  amountOfCakePerYear= multiplier * cakePerYear;
  totalValueOfCake=amountOfCakePerYear*priceOfCake;
  farmBaseReward= BigInt(totalValueOfCake)/liquidity*BigInt(100); //Percentage
  totalFee= (volume24h*0.0017)*365;
  lpReward = BigInt(totalFee)/liquidity*BigInt(100);
  APR= farmBaseReward + lpReward;
    // ListPar(i,result[1],"","",APR,liquidity,multiplier);
    document.getElementById(`APR_${idPool}`).innerHTML =APR;
    
  
}



//Load Data in smart contract when open browser 
var loadData = async(result2) => {
  //Load list of Farm Pools
  // var result =[]
  for(var i=0 ;i<(result2.length);i++){
  // result = result2.slice(i*4,4+i*4);
  // var liquidity,farmBaseReward,totalFee,volume24h;
  // var priceOfCake = 20;//getAmountOut in router????
  // cakePerYear=cakePerblock*60*60*24*365;
  // liquidity =result[2];
  // multiplier = result[0];
  // volume24h = 1000// dang fix cung theo usd
  // totalValueOfCake= cakePerYear *priceOfCake;
  // amountOfCakePerYear= multiplier * cakePerYear;
  // totalValueOfCake=amountOfCakePerYear*priceOfCake;
  // farmBaseReward= totalValueOfCake/liquidity*100; //Percentage
  // totalFee= (volume24h*0.0017)*365;
  // lpReward = totalFee/liquidity*100;
  // APR= farmBaseReward + lpReward;
  // APR = result[3];
    ListPar(i,result[1],"","",APR,liquidity,multiplier);
    document.getElementById(`liquidity_${i}`).value =liquidity;
    document.getElementById(`APR_${i}`).value =APR;
    document.getElementById(`multiplier_${i}`).value =multiplier;
    
  }
}
var loadDataCompound = ()=>{
  fetch('/compoundToken')
  .then(function(response) {
    if (!response.ok) {
      throw Error(response.statusText);
    }
    // Read the response as json.
    return response.json();
  })
  .then(function(responseAsJson) {
    // Do stuff with the JSON
    var html1 = responseAsJson.map(function(token){  
    ListSupply(token.address,token.supplyapr,"","")
    ListBorrow(token.address,token.borrowapr,"",token.liquidity)
    });
    for (var key in responseAsJson){
      getInfo(responseAsJson[key].address);
      console.log("responseAsJson[key].address:",responseAsJson[key].address)
    }
    
  })
  .catch(function(error) {
    console.log('Looks like there was a problem: \n', error);
  });
    
}
var getInfo = (tokenAdd)=>{
  fetch('/compoundUser')
  .then(function(response) {
    if (!response.ok) {
      throw Error(response.statusText);
    }
    // Read the response as json.
    return response.json();
  })
  .then(function(responseAsJson) {
    var html1 = responseAsJson[0].supplybalance
    var html2 = responseAsJson[0].borrowbalance
      for (var key in html1){
        if(key==tokenAdd){
        document.getElementById(`idSupplyWallet_${key}`).innerHTML =html1[key];
          console.log("html1[key]:",html1[key])
        }
      }

      for (var key in html2){
        if(key==tokenAdd){
          document.getElementById(`idBorWallet_${key}`).innerHTML =html2[key];
          console.log("html2[key]:",html2[key])
        }
      }
  })
  .catch(function(error) {
    console.log('Looks like there was a problem: \n', error);
  }); 
}
//List Supply CompoundTab Table
var ListSupply=async(a,b,c,d)=>{
  var listPar = ` <tr >
  <th class="text-center" width="25%" id="idSupplyAsset_${a}">${a}</th>
  <th class="text-center" width="25%" id="idSupplyApy_${a}">${b}</th>
  <th class="text-center" width="25%" id="idSupplyWallet_${a}">${c}</th>
  <th class="text-center" width="25%" id="idSupplyColl_${a}">${d}</th>
  </tr>`
$('.supply_table tbody').append(listPar)   

}
//List Borrow CompoundTab Table
var ListBorrow=async(a,b,c,d)=>{
  var listPar = ` <tr >
  <th class="text-center" width="25%" id="idBorAsset_${a}">${a}</th>
  <th class="text-center" width="25%" id="idBorApy_${a}">${b}</th>
  <th class="text-center" width="25%" id="idBorWallet_${a}">${c}</th>
  <th class="text-center" width="25%" id="idBorLiquidity_${a}">${d}</th>
  </tr>`
$('.borrow_table tbody').append(listPar)   

}


//List Participant Table
var ListPar=async(a,b,c,d,e,f,g)=>{
  var listPar = ` <tr >
  <th class="text-center" width="15%" id="idPool_${a}">${a}</th>
  <th class="text-center" width="15%" id="name_${a}">${b}</th>
  <th class="text-center" width="15%" id="boost_${a}">${c}</th>
  <th class="text-center" width="15%" id="earn_${a}">${d}</th>
  <th class="text-center" width="25%" id="APR_${a}" value=0>${e} </th>
  <th class="text-center" width="15%" id="liquidity_${a}" value=0>${f}</th>
  <th class="text-center" width="15%" id="multiplier_${a}" value=0>${g}</th>

</tr>`
$('.participants_table tbody').append(listPar)   

}

// getElemenetById("farm-table").addEventListener("click",handleFarmTable);

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
getElemenetById("compoundTab").addEventListener("click",getCompoundTab);
getElemenetById("stakingTab").addEventListener("click",getStakingTab);

var sendMessage = (msg) => {
  console.log(msg);
  socket.send(JSON.stringify(msg));
};
// getElemenetById("reset-qr-btn").addEventListener("click", handleResetQR);
// var handleResetQR = () => {
//   eraseAvailableQR();
// };
//addLiquidityMTD value=10000
// f305d719
// 000000000000000000000000b9C40c5054333975e4fEE5b2972f2481422CD48D token
// 0000000000000000000000000000000000000000000000000000000000002710 amountdesired
// 0000000000000000000000000000000000000000000000000000000000000001 amountTokenMin
// 0000000000000000000000000000000000000000000000000000000000000001 amountMTDMin
// 000000000000000000000000d85ae9a6ef6185aea70b1b18c3d3bfd1253ea74e to
// 0000000000000000000000000000000000000000000000056bc75e2d63100000 time

// 38ed1739
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
// 000000000000000000000000000000000000000000000000000000000010f447
// 0000000000000000000000000000000000000000000000000000000000000000
// 00000000000000000000000000000000000000000000000000000000000000a0
// 000000000000000000000000f81751083e57b4ab07929a6ec931cca20464f393
// 000000000000000000000000000000000000000000000001007a33ee6d770000
// 0000000000000000000000000000000000000000000000000000000000000002
// 00000000000000000000000045b1b47617195012803515df966009AFB708Ad25
// 0000000000000000000000006B5B7Adc611AB1C9434fdf5A3357b8daDDA0b703

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
