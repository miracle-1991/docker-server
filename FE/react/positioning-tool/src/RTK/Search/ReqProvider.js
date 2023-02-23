import React, {createContext, useContext, useEffect, useState} from "react";

//创建search组件公用的context
const ReqContext = createContext();
export const useReq = () => useContext(ReqContext);

//整个search组件都将是该组件的子组件
export default function ReqProvider({children}) {
    const [req, setReq] = useState({
        rtk: {
            pmode: [
                "kinema",
                "static"
            ],
            navsys: [
                "gps-gal-glo",
                "gps-gal-glo-cmp"
            ],
            obsnav: {
                roverobs: "/data/20230215/RoadTesting20230215/forestrouteloop2/Pixel6/gnss_log_2023_02_15_14_44_58.23o",
                stationobs: "/data/20230215/RoadTesting20230215/SIN100SGP_S_20230460000_01D_30S_MO.rnx",
                stationnav: "/data/20230215/RoadTesting20230215/BRDC00IGS_R_20230460000_01D_MN.rnx"
            }
        },
        outputpath: "/data"
    });

    const pmodelist = [{name: "dgps"}, {name: "kinema"}, {name: "static"}]
    const navsyslist = [{name: "gps"},{name: "gal"},{name: "glo"},{name:"cmp"}]
    const [pmodeCheckState, setPmodeCheckState] = useState(new Array(pmodelist.length).fill(false));
    const [navsysCheckState, setNavSysCheckState] = useState(new Array(navsyslist.length).fill(false));
    const [rtkuri, setRtkUri] = useState("http://localhost:8001/DEMO5/RTK/demo5")
    const [rtkData, setRtkData] = useState();
    const [error, setError] = useState();
    const [rtkloading, setRtkLoading] = useState(true);
    const [rtkProcessingState, setRtkLogProcessingState] = useState(0.0);

    const onPmodeCheckChange = (position) => {
        const updatedState = pmodeCheckState.map((item, index) => (
            index === position ? ! item : item
        ));
        setPmodeCheckState(updatedState);
        let plist = [];
        for (let i = 0; i < updatedState.length; i++) {
            if (updatedState[i] === true) {
                plist.push(pmodelist[i].name)
            }
        }
        let newReq = {...req}
        newReq.rtk.pmode = plist
        setReq(newReq)
        console.log("onPmodeCheckChange: " + plist + "   " + JSON.stringify(req))
    }

    const onNavSysCheckChange = (position) => {
        const updatedState = navsysCheckState.map((item, index) => (
            index === position ? !item : item
        ));
        setNavSysCheckState(updatedState);
        let nlist = [];
        for (let i = 0; i < updatedState.length; i++) {
            if (updatedState[i] === true) {
                nlist.push(navsyslist[i].name)
            }
        }
        let newReq = {...req}
        newReq.rtk.navsys = [ nlist.join("-") ];
        setReq(newReq);
        console.log("onNavSysCheckChange: " + nlist + " " + JSON.stringify(req))
    }

    const onDataSuccess = data => {
        console.log(data)
    }

    const onRTKDataLoading = () => {
        setRtkData();
        let reqCnt = 0;
        let curprocessingState = 0
        let timerId = setInterval(() => {
            reqCnt++
            if (curprocessingState >= 1 || reqCnt === 100) {
                clearInterval(timerId)
            }
            fetch("http://localhost:8001/DEMO5/RTK/rtkProcessing")
                .then(rsp => rsp.json())
                .then(rsp => {
                    curprocessingState = rsp.processing
                    setRtkLogProcessingState(curprocessingState)
                }).catch(error => console.error(error))
        }, 1000);
    }

    useEffect(() => {
        console.log("rtkProcessingState = " + rtkProcessingState)
    }, [rtkProcessingState])

    useEffect(() => {
        if (error != null) {
            alert(error)
        }
    }, [error])

    useEffect(() => {
        console.log(rtkData)
    }, [rtkData])

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
    const onOutputPathChange = pathStr => {
        setReq({
            ...req,
            outputpath: pathStr
        })
        console.log("onOutputPathChange: " + pathStr + "   " + JSON.stringify(req))
    }

    const onRTKRoverObsPathChange = pathstr => {
        let newReq = {...req}
        newReq.rtk.obsnav.roverobs = pathstr;
        setReq(newReq);
        console.log("onRTKRoverObsPathChange: " + pathstr + "   " + JSON.stringify(req))
    }

    const onRTKStationObsPathChange = pathstr => {
        let newReq = {...req}
        newReq.rtk.obsnav.stationobs = pathstr;
        setReq(newReq);
        console.log("onRTKStationObsPathChange: " + pathstr + "   " + JSON.stringify(req))
    }

    const onRTKStationNavPathChange = pathstr => {
        let newReq = {...req}
        newReq.rtk.obsnav.stationnav = pathstr;
        setReq(newReq);
        console.log("onRTKStationNavPathChange: " + pathstr + "   " + JSON.stringify(req))
    }

    return (
        <ReqContext.Provider value={{
            req, onStartTimeChange, onEndTimeChange, onfilterChange, onOutputPathChange,
            rtkuri, setRtkUri,
            rtkData, setRtkData,
            error, setError,
            rtkloading, setRtkLoading,rtkProcessingState,
            onRTKDataLoading, onDataSuccess, onDataError,
            pmodelist, pmodeCheckState, onPmodeCheckChange,
            navsyslist, navsysCheckState, onNavSysCheckChange,
            onRTKRoverObsPathChange,onRTKStationObsPathChange,onRTKStationNavPathChange
        }}>
            {children}
        </ReqContext.Provider>
    )
}