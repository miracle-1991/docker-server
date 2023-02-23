import React from "react";

export default function PmodeItem({name, index, checked, onchange}) {
    return (
        <label className="checkbox-inline">
            <input type="checkbox" id={`pmode-checkbox-${index}`} checked={checked[index]} onChange={()=> onchange(index)}/> {name}
        </label>
    )
}