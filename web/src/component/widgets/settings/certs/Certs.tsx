import {useState} from "react";

import {useRouterCertList} from "../../../../api/cert/hook";
import {LinearProgressStateful} from "../../../view/progress/LinearProgressStateful";
import {MenuWrapper} from "../menu/MenuWrapper";
import {CertsList} from "./CertsList";
import {CertsNew} from "./CertsNew";
import {CertsTab, TABS} from "./CertsTab";

export function Certs() {
    const [tab, setTab] = useState(0)
    const type = TABS[tab].type
    const query = useRouterCertList(type)
    const {data, error, isFetching} = query

    return (
        <MenuWrapper>
            <CertsTab tab={tab} setTab={setTab}/>
            <CertsNew type={type}/>
            <LinearProgressStateful color={"inherit"} loading={isFetching} line/>
            <CertsList certs={data} error={error}/>
        </MenuWrapper>
    )
}
