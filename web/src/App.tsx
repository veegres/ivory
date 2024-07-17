import {Box} from "@mui/material";
import {Header} from "./component/section/Header";
import {Body} from "./component/section/Body";
import {Footer} from "./component/section/Footer";
import {AuthType, SxPropsMap} from "./type/general";
import {useRouterInfo} from "./router/general";
import {ReactNode} from "react";

const SX: SxPropsMap = {
    // we need -8px for bottom scroll, right scroll is always shown that is why we don't need calc
    box: {display: "flex", flexDirection: "column", gap: 2, flexWrap: "nowrap", height: "calc(100vh - 8px)", minWidth: "1010px"},
    body: {display: "flex", alignItems: "center", justifyContent: "center", flexGrow: 1},
}

type Props = {
    children?: ReactNode,
}

export function App(props: Props) {
    const info = useRouterInfo()
    const show = !info.isPending && !!info.data && info.data.secret.key && info.data.auth.authorised

    return (
        <Box sx={SX.box}>
            <Box>
                <Header company={info.data?.company ?? "VEEGRES"} show={!props.children && show} auth={info.data?.auth.type ?? AuthType.NONE}/>
            </Box>
            <Box sx={SX.body}>
                {props.children ? props.children : <Body info={info}/>}
            </Box>
            <Box>
                <Footer label={info.data?.version.label ?? "none"} tag={info.data?.version.tag ?? "master"}/>
            </Box>
        </Box>
    );
}
