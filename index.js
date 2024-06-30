let data = []
let input = document.getElementById("json-file-input");
input.onchange = event => {
    let files = event.target.files;
    let waitCounter = 0;
    for(let i = 0; i < files.length; i++) {
        let reader = new FileReader();
        reader.readAsText(files[i],'UTF-8');  
        
        reader.onload = readerEvent => {
            let rawdata = JSON.parse(readerEvent.target.result);
    
            
            for(let key in rawdata) {
                let item = rawdata[key]
                if (data[item.Name] == null) {
                    data[item.Name] = []
                    data[item.Name].layout = {
                        title: item.Name,
                        xaxis: {
                            tickangle: -30
                          },
                          yaxis: {
                            zeroline: false,
                            gridwidth: 2,
                            title: "Computation time, ms"
                          },
                          bargap :0.5,
                          height: 500,
                          width: 700,
                          showlegend: true
                    }
                }
                if (data[item.Name].blockSizes == null) {
                    data[item.Name].blockSizes = []
                }

                if(!data[item.Name].blockSizes.find((element) => element.blockSize == item.BlockSize)) {
                    data[item.Name].blockSizes.push({
                        blockSize: item.BlockSize,
                        x:[],
                        y:[],
                        type: 'bar',
                        name: "Block Size " + item.BlockSize.toString()
                    })
                }
               
                for(bs in data[item.Name].blockSizes) {
                    bss = data[item.Name].blockSizes[bs]
                    if (bss.blockSize == item.BlockSize) {
                        bss.x.push("Matrix Size: "+item.MatrixSize.toString())
                        bss.y.push(item.AvgTimeElapsed / 1000000)
                    }
                }
            }

            for(let name in data) {
                let newBlockSizesArray = [];
                for (let i = 0; i < data[name].blockSizes.length; i++) {
                    if (data[name].blockSizes[i].blockSize == 0) {
                        for(let j = 0; j < data[name].blockSizes[i].x.length; j++) {
                            newBlockSizesArray.push({
                                blockSize: -1,
                                x:[data[name].blockSizes[i].x[j]],
                                y:[data[name].blockSizes[i].y[j]],
                                type: 'bar',
                                name: "Iteration#" + (j+1).toString()
                            });
                        }

                    }
                }
                if (newBlockSizesArray.length > 0) {
                    data[name].blockSizes = newBlockSizesArray;
                }
            }

            waitCounter++;

            if(waitCounter == files.length) {
                for(let name in data) {
                    let graphDiv = document.createElement("div")
                    graphDiv.id = "plotDiv-"+name
                    document.body.appendChild(graphDiv);
                    console.log(data[name])
                    Plotly.newPlot(graphDiv.id, data[name].blockSizes, data[name].layout);
        
                    let bargap = document.createElement("hr")
                    document.body.appendChild(bargap);
                    bargap = document.createElement("hr")
                    document.body.appendChild(bargap);
                }
            }
         }
    }

    
}
