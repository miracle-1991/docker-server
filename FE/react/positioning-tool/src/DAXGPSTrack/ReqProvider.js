import React, {createContext, useContext, useEffect, useState} from "react";

//创建search组件公用的context
const ReqContext = createContext();
export const useReq = () => useContext(ReqContext);

//整个search组件都将是该组件的子组件
export default function ReqProvider({children}) {
    const [req, setReq] = useState({
        driverid: 25463,
        time: "2023-02-09 23:57:00",
        forward: 10,
        backward: 10,
        outputpath: "/data"
    });

    const [uri, setUri] = useState("http://localhost:8003/driver")
    const [data, setData] = useState();
    const [error, setError] = useState();
    const [loading, setLoading] = useState(true);

    const onDataSuccess = data => {
        console.log(data)
    }

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

    const onDriveridChange = tint => {
        setReq({
            ...req,
            driverid: tint
        })
        console.log("onDriveridChange: " + tint + "   " + JSON.stringify(req))
    }
    const onTimeChange = tstr => {
        setReq({
            ...req,
            time: tstr
        })
        console.log("onTimeChange: " + tstr + "   " + JSON.stringify(req))
    }
    const onForwardChange = tint => {
        setReq({
            ...req,
            forward: tint
        })
        console.log("onForwardChange: " + tint + "   " + JSON.stringify(req))
    }
    const onBackwardChange = tint => {
        setReq({
            ...req,
            backward: tint
        })
        console.log("onBackwardChange: " + tint + "   " + JSON.stringify(req))
    }
    const onOutputpathChange = tstr => {
        setReq({
            ...req,
            outputpath: tstr
        })
        console.log("onOutputpathChange: " + tstr + "   " + JSON.stringify(req))
    }
    return (
        <ReqContext.Provider value={{
            req, onDriveridChange, onTimeChange, onForwardChange, onBackwardChange,onOutputpathChange,
            uri, setUri, data, setData, error, setError, loading, setLoading,
            onDataSuccess, onDataError
        }}>
            {children}
        </ReqContext.Provider>
    )
}