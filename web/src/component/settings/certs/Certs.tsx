import {useState} from "react";
import {CertsNew} from "./CertsNew";
import {CertsList} from "./CertsList";
import {CertType} from "../../../app/types";
import {MenuWrapper} from "../menu/MenuWrapper";
import {CertsTab, TABS} from "./CertsTab";

export type CertTypeProps = {
    type: CertType,
}

export function Certs() {
    const [tab, setTab] = useState(0)

    return (
        <MenuWrapper>
            <CertsTab tab={tab} setTab={setTab}/>
            <CertsNew type={TABS[tab].type}/>
            <CertsList type={TABS[tab].type}/>
        </MenuWrapper>
    )
}
