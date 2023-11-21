import {useState} from "react";
import {CertsNew} from "./CertsNew";
import {CertsList} from "./CertsList";
import {MenuWrapper} from "../menu/MenuWrapper";
import {CertsTab, TABS} from "./CertsTab";
import {useQuery} from "@tanstack/react-query";
import {certApi} from "../../../app/api";
import {LinearProgressStateful} from "../../view/progress/LinearProgressStateful";

export function Certs() {
    const [tab, setTab] = useState(0)
    const type = TABS[tab].type
    const query = useQuery({queryKey: ["certs", type], queryFn: () => certApi.list(type)})
    const {data, error, isFetching} = query

    return (
        <MenuWrapper>
            <CertsTab tab={tab} setTab={setTab}/>
            <CertsNew type={type}/>
            <LinearProgressStateful color={"inherit"} isFetching={isFetching} line/>
            <CertsList certs={data} error={error}/>
        </MenuWrapper>
    )
}
