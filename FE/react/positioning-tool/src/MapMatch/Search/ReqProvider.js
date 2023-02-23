import React, {createContext, useContext, useEffect, useState} from "react";

//创建search组件公用的context
const ReqContext = createContext();
export const useReq = () => useContext(ReqContext);

//整个search组件都将是该组件的子组件
export default function ReqProvider({children}) {
    const [req, setReq] = useState({
        filepath: "/data/result/forestrouteloop1/Note20Ultra",
        filename: "rtkFilter-driverID-13770990-adr.csv",
        latitude_column_name: "lat",
        longitude_column_name: "lon",
        timestamp_column_name: "timestamp"
    });

    const [snapuri, setSnapUri] = useState("http://localhost:8000/snap")
    const [snapData, setSnapData] = useState();
    const [error, setError] = useState();
    const [snaploading, setSnapLoading] = useState(true);

    const onDataSuccess = data => {
        console.log(data)
    }

    useEffect(() => {
        if (error != null) {
            alert(error)
        }
    }, [error])

    useEffect(() => {
        console.log(snapData)
    }, [snapData])

    const onDataError = error => {
        if (error != null) {
            console.error(error)
        }
    }

    const onFilepathChange = tstr => {
        setReq({
            ...req,
            filepath: tstr
        })
        console.log("onFilepathChange: " + tstr + "   " + JSON.stringify(req))
    }
    const onFilenameChange = tstr => {
        setReq({
            ...req,
            filename: tstr
        })
        console.log("onFilenameChange: " + tstr + "   " + JSON.stringify(req))
    }
    const onLatitudeColumnNameChange = tstr => {
        setReq({
            ...req,
            latitude_column_name: tstr
        })
        console.log("onLatitudeColumnNameChange: " + tstr + "   " + JSON.stringify(req))
    }
    const onLongitudeColumnNameChange = tstr => {
        setReq({
            ...req,
            longitude_column_name: tstr
        })
        console.log("onLongitudeColumnNameChange: " + tstr + "   " + JSON.stringify(req))
    }

    const onTimestampColumnNameChange = tstr => {
        setReq({
            ...req,
            timestamp_column_name: tstr
        })
        console.log("onTimestampColumnNameChange: " + tstr + "   " + JSON.stringify(req))
    }

    return (
        <ReqContext.Provider value={{
            req, error, setError, snapuri, setSnapUri, snapData, setSnapData, snaploading, setSnapLoading,
            onDataError, onDataSuccess, onFilepathChange, onFilenameChange, onLatitudeColumnNameChange,
            onLongitudeColumnNameChange, onTimestampColumnNameChange
        }}>
            {children}
        </ReqContext.Provider>
    )
}