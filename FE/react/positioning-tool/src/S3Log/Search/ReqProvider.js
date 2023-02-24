import React, {createContext, useContext, useEffect, useState} from "react";

//创建search组件公用的context
const ReqContext = createContext();
export const useReq = () => useContext(ReqContext);

//整个search组件都将是该组件的子组件
export default function ReqProvider({children}) {
    const [req, setReq] = useState({
        starttime: "2023-02-15 07:54:10",
        endtime: "2023-02-15 08:21:19",
        filter: ["rtkFilter","driverID:13773457"],
        prefix: "location-engine",
        outputpath: "/data"
    });

    const [uri, setUri] = useState("http://" + process.env.REACT_APP_LOCALHOST + ":8000/downlog")
    const [data, setData] = useState();
    const [error, setError] = useState();
    const [loading, setLoading] = useState(true);
    const [processingState, setProcessingState] = useState(0.0);

    const onDataSuccess = data => {
        console.log(data)
    }

    const onDataLoading = () => {
        setData();
        let reqCnt = 0;
        let curprocessingState = 0
        let timerId = setInterval(() => {
            reqCnt++
            if (curprocessingState >= 1 || reqCnt === 300) {
                clearInterval(timerId)
            }
            fetch("http://" + process.env.REACT_APP_LOCALHOST +":8000/downlogProcessing")
                .then(rsp => rsp.json())
                .then(rsp => {
                    curprocessingState = rsp.processing
                    setProcessingState(curprocessingState)
                }).catch(error => console.error(error))
        }, 1000);
    }

    useEffect(() => {
        console.log("processingState = " + processingState)
    }, [processingState])

    useEffect(() => {
        if (error != null) {
            alert(error)
        }
    }, [error])

    useEffect(() => {
        console.log(data)
    }, [data])

    const onDataError = error => {
        if (error != null) {
            console.error(error)
        }
    }

    const onStartTimeChange = tstr => {
        setReq({
            ...req,
            starttime: tstr
        })
        console.log("onStartTimeChange: " + tstr + "   " + JSON.stringify(req))
    }
    const onEndTimeChange = tstr => {
        setReq({
            ...req,
            endtime: tstr
        })
        console.log("onEndTimeChange: " + tstr + "   " + JSON.stringify(req))
    }
    const onfilterChange = filterliststr => {
        let flit = filterliststr.split(",")
        setReq({
            ...req,
            filter: flit
        })
        console.log("onEndTimeChange: " + filterliststr + "   " + JSON.stringify(req))
    }

    const onPrefixChange = tstr => {
        setReq({
            ...req,
            prefix: tstr
        })
        console.log("onPrefixChange: " + tstr + "  " + JSON.stringify(req))
    }

    const onOutputPathChange = pathStr => {
        setReq({
            ...req,
            outputpath: pathStr
        })
        console.log("onOutputPathChange: " + pathStr + "   " + JSON.stringify(req))
    }
    return (
        <ReqContext.Provider value={{
            req, onStartTimeChange, onEndTimeChange, onfilterChange,onPrefixChange, onOutputPathChange,
            uri, setUri, data, setData, error, setError, loading, setLoading,
            onDataSuccess, onDataLoading, onDataError, processingState
        }}>
            {children}
        </ReqContext.Provider>
    )
}