import React from "react";

export default function NavSysItem({name, index, checked, onchange}) {
    return (
        <label className="checkbox-inline">
            <input type="checkbox" id={`navsys-checkbox-${index}`} checked={checked[index]} onChange={()=> onchange(index)}/> {name}
        </label>
    )
}