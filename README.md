# AtomicOTCswap
#### Atomic Swap tool bitcoin &amp; altcoins without the need of full nodes on-chain
###### This atomic swap tool is heavily based on Decred's atomic swap tool: https://github.com/decred/atomicswap

##### Supported assets
- [ ] Bitcoin ([Bitcoin Core](https://github.com/bitcoin/bitcoin))
- [ ] Decred ([dcrwallet](https://github.com/decred/dcrwallet))
- [x] Viacoin ([Viacoin Core](https://github.com/viacoin/viacoin))
- [x] Litecoin ([Litecoin Core](https://github.com/litecoin-project/litecoin))
- [ ] Vertcoin ([Vertcoin Core](https://github.com/vertcoin/vertcoin))
- [ ] Zcoin ([Zcoin Core](https://github.com/zcoinofficial/zcoin))
- [ ] Bitcoin Testnet ([Bitcoin Core](https://github.com/bitcoin/bitcoin))
- [ ] Decred Tesnet ([dcrwallet](https://github.com/decred/dcrwallet))
- [ ] Litecoin Tesnet ([Litecoin Core](https://github.com/litecoin-project/litecoin))

###### Disclaimer
This swap tool and the UI is highly experimental. <br>
One could easily have his funds "stuck". Only use with small amounts, this tool is not for "the average Joe" for now.<br><br>
The author takes no responsibility
<br><br>
This project is not finished, it lacks unit testing and this tool written by a Golang newbie :) <br>
I still need to do a lot of refactoring & cleaning up the code base. It's far from finished but it works :) </br>

Feel free to create pull requests and help improving. 

![alt text](https://github.com/romanornr/AtomicOTCswap/blob/master/screenshots/1.png?raw=true)


###### API
The atomic swap tool comes with an API so others can build on top of it or for example host the atomic swap server and let users connect externily and preferably connect to their own host.

One could also host the atomic swap server as a site. (Users may have to "trust" which defeats the purpose a bit but it's doable)

API documentation will be added later.


###### Minimum Recommended Specifications

- Go 1.10 or 1.11
* Linux


  Installation instructions can be found here: https://golang.org/doc/install.
  It is recommended to add `$GOPATH/bin` to your `PATH` at this point.

###### setup
``cd ~/go/src/github.com/``

``git clone git@gitlab.com:romanornr/atomicotcswap``

``cd atomicotcswap``

``dep ensure`` 


dep is a dependency management tool for Go. It requires Go 1.9 or newer to compile.
https://github.com/golang/dep



Credits
=======
  
  - Decred Project
    * [jrick](https://github.com/jrick)<br/>
      [Author Decred atomic swap](https://github.com)<br>
      
    * [dajohi](https://github.com/dajohi)<br/>
      [Author Decred atomic swap](https://github.com/decred/atomicswap)<br>
      <br>Special credits to the Decred team. This tool is heavily based on Decred's atomic swap tool
      
  - Frontend
    * [rockstardev](https://github.com/rockstardev)
