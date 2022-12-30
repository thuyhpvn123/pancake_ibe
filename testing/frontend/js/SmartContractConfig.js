// const SMAToken = {
//   address: "b9C40c5054333975e4fEE5b2972f2481422CD48D",
//   approve: {
//     name: "approve",
//     parameter: [
//       //Hash
//       {
//         type: "hash",
//         value: "095ea7b3",
//       },

//       //Parameter
//       {
//         name: "spender",
//         type: "address",
//         value: "",
//       },
//       {
//         name: "amount",
//         type: "num",
//         value: "",
//       },
//     ],
//   },
// };

// const SMBToken = {
//   address: "148b3c21920A743625974Bf7E7f3b8D675534b74",
//   approve: {
//     name: "approve",
//     parameter: [
//       //Hash
//       {
//         type: "hash",
//         value: "095ea7b3",
//       },

//       //Parameter
//       {
//         name: "spender",
//         type: "address",
//         value: "",
//       },
//       {
//         name: "amount",
//         type: "num",
//         value: "",
//       },
//     ],
//   },
// };
const SMToken = {
  // address: "148b3c21920A743625974Bf7E7f3b8D675534b74",
  approve: {
    name: "approve",
    parameter: [
      //Hash
      {
        type: "hash",
        value: "095ea7b3",
      },

      //Parameter
      {
        name: "spender",
        type: "address",
        value: "",
      },
      {
        name: "amount",
        type: "num",
        value: "",
      },
    ],
  },
};

const LPToken = {
  address: "79c0494029712423ae93a8a7711fca15721b1929",
  approve: {
    name: "approve",
    parameter: [
      //Hash
      {
        type: "hash",
        value: "095ea7b3",
      },

      //Parameter
      {
        name: "spender",
        type: "address",
        value: "",
      },
      {
        name: "amount",
        type: "num",
        value: "",
      },
    ],
  },
};

const MTDToken = {
  address: "2a955DBfa8ac000584Fd5c006af60B23F0888F6e",
};

const SMRouter = {
  address: "2b00fAE7F567fEf7b1dfCF64067D344C8B055b0c",
  addLiquidity: {
    name: "addLiquidity",
    parameter: [
      //Hash
      {
        type: "hash",
        value: "e8e33700",
      },

      //Parameter
      {
        name: "tokenA",
        type: "address",
        value: "",
      },
      {
        name: "tokenB",
        type: "address",
        value: "",
      },
      {
        name: "amountADesiredA",
        type: "num",
        value: "",
      },
      {
        name: "amountADesiredB",
        type: "num",
        value: "",
      },
      {
        name: "amountAMin",
        type: "num",
        value: 0,
      },
      {
        name: "amountBMin",
        type: "num",
        value: 0,
      },
      {
        name: "to",
        type: "address",
        value: "",
      },
      {
        name: "deadline",
        type: "num",
        value: "",
      },
    ],
  },
  addLiquidityMTD: {
    name: "addLiquidityMTD",
    parameter: [
      //Hash
      {
        type: "hash",
        value: "f305d719",
      },

      //Parameter
      {
        name: "tokenA",
        type: "address",
        value: "",
      },
      {
        name: "amountADesiredA",
        type: "num",
        value: "",
      },
      {
        name: "amountAMin",
        type: "num",
        value: 0,
      },
      {
        name: "amountBMin",
        type: "num",
        value: 0,
      },
      {
        name: "to",
        type: "address",
        value: "",
      },
      {
        name: "deadline",
        type: "num",
        value: "",
      },
    ],
  },

  removeLiquidity: {
  name: "removeLiquidity",
  parameter: [
    //Hash
    {
      type: "hash",
      value: "baa2abde",
    },

    //Parameter
    {
      name: "tokenA",    //1
      type: "address",
      value: "",
    },
    {
      name: "tokenB",    //2
      type: "address",
      value: "",
    },
    {
      name: "liquidity",   //3
      type: "num",
      value: "",
    },
    {
      name: "amountAMin",
      type: "num",
      value: 1,
    },
    {
      name: "amountBMin",
      type: "num",
      value: 1,
    },
    {
      name: "to",      //6
      type: "address",
      value: "",
    },
    {
      name: "deadline", //7
      type: "num",
      value: "",
    },
  ],
  },


  swap: {
    name: "swap",
    parameter: [
      //Hash
      {
        type: "hash",
        value: "38ed1739",
      },

      //Parameter
      {
        name: "amountADesired",   //1
        type: "num",
        value: "",
      },
      {
        name: "tokenBmin",         //2
        type: "num",
        value: 0,
      },
      {
        name: "offsetArray",  //3
        type: "num",
        value: 160,
      },
      {
        name: "to",   //4
        type: "address",
        value: "",
      },
      {
        name: "deadline",    //5
        type: "num",
        value: "",
      },
      {
        name: "arrayLength",   //6
        type: "num",
        value: 2,
      },
      {
        name: "tokenA",     //7
        type: "address",
        value: "",
      },
      {
        name: "tokenB",   //8
        type: "address",
        value: "",
      },
    ],
  },
  swapMTD: {   //swapExactETHForTokens
    name: "swapMTD",
    parameter: [
      //Hash
      {
        type: "hash",
        value: "7ff36ab5",
      },

      //Parameter
      {
        name: "amountOutMin",   //1
        type: "num",
        value: 0,
      },
      // {
      //   name: "tokenMTD",         //
      //   type: "num",
      //   value: "",
      // },
      {
        name: "offsetArray",  //2
        type: "num",
        value: 128,
      },
      {
        name: "to",   //3
        type: "address",
        value: "",
      },
      {
        name: "deadline",    //4
        type: "num",
        value: "",
      },
      {
        name: "arrayLength",   //5
        type: "num",
        value: 2,
      },
      {
        name: "tokenA",     //6
        type: "address",
        value: "",
      },
      {
        name: "tokenB",   //7
        type: "address",
        value: "",
      },
    ],
  },
  swapToMTD: {
    name: "swap",
    parameter: [
      //Hash
      {
        type: "hash",
        value: "18cbafe5",
      },

      //Parameter
      {
        name: "amountADesired",   //1
        type: "num",
        value: "",
      },
      {
        name: "tokenOutMin",         //2
        type: "num",
        value: 0,
      },
      {
        name: "offsetArray",  //3
        type: "num",
        value: 160,
      },
      {
        name: "to",   //4
        type: "address",
        value: "",
      },
      {
        name: "deadline",    //5
        type: "num",
        value: "",
      },
      {
        name: "arrayLength",   //6
        type: "num",
        value: 2,
      },
      {
        name: "tokenA",     //7
        type: "address",
        value: "",
      },
      {
        name: "tokenB",   //8
        type: "address",
        value: "",
      },
    ],
  },


};


const DEADLINE = 18481141120000000000;
