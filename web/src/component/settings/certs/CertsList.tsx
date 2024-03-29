import {Collapse} from "@mui/material";
import {TransitionGroup} from "react-transition-group";
import {CertsItem} from "./CertsItem";
import {Cert, CertMap} from "../../../type/cert";
import {InfoAlert} from "../../view/box/InfoAlert";
import {ErrorSmart} from "../../view/box/ErrorSmart";

type Props = {
    certs?: CertMap,
    error: any,
}

export function CertsList(props: Props) {
    const {certs, error} = props
    if (error) return <ErrorSmart error={error}/>

    const list = Object.entries<Cert>(certs ?? {})
    if (list.length === 0) return <InfoAlert text={"There is no files yet"}/>

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
