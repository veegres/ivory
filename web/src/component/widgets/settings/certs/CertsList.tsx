import {Collapse} from "@mui/material"
import {TransitionGroup} from "react-transition-group"

import {Cert, CertMap} from "../../../../api/cert/type"
import {ErrorSmart} from "../../../view/box/ErrorSmart"
import {InfoAlert} from "../../../view/box/InfoAlert"
import {CertsItem} from "./CertsItem"

type Props = {
    certs?: CertMap,
    error: any,
}

export function CertsList(props: Props) {
    const {certs, error} = props
    if (error) return <ErrorSmart error={error}/>

    const list = Object.entries<Cert>(certs ?? {})
    if (list.length === 0) return <InfoAlert text={"There are no files yet"}/>

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
