import {Collapse} from "@mui/material"
import {TransitionGroup} from "react-transition-group"

import {ErrorSmart} from "../../../shared/component/box/ErrorSmart"
import {NoBox} from "../../../shared/component/box/NoBox"
import {Cert, CertMap} from "../api/CertType"
import {CertsItem} from "./CertsItem"

type Props = {
    certs?: CertMap,
    error: any,
}

export function CertsList(props: Props) {
    const {certs, error} = props
    if (error) return <ErrorSmart error={error}/>

    const list = Object.entries<Cert>(certs ?? {})
    if (list.length === 0) return <NoBox text={"There are no cert files yet"}/>

    return (
        <TransitionGroup appear={false}>
            {list.map(([key, cert]) => (
                <Collapse key={key}>
                    <CertsItem uuid={key} cert={cert}/>
                </Collapse>
            ))}
        </TransitionGroup>
    )
}
