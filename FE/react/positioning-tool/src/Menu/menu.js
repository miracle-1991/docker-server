import React from "react";
import "../dist/css/zui.css"

export default function Menu({choose = "README"}) {
    return (
        <nav className="navbar navbar-default" role="navigation">
            <ul className="nav navbar-nav nav-justified">
                <li className={choose === "driver" ? "active": null}><a href="/">DAX-GPS</a></li>
                <li className={choose === "booking" ? "active": null}><a href="/booking">BOOKING-GPS</a></li>
                <li className={choose === "s3log" ? "active": null}><a href="/s3log">S3LOG</a></li>
                <li className={choose === "snap" ? "active": null}><a href="/snap">SNAP</a></li>
                <li className={choose === "rtk" ? "active": null}><a href="/rtk">RTK</a></li>
            </ul>
        </nav>
    )
}
