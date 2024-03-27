import {Box} from "@mui/material";
import {Header} from "./component/section/Header";
import {Body} from "./component/section/Body";
import {useQuery} from "@tanstack/react-query";
import {GeneralApi} from "./app/api";
import {Footer} from "./component/section/Footer";
import {AuthType, SxPropsMap} from "./type/general";

const SX: SxPropsMap = {
    // we need -8px for bottom scroll, right scroll is always shown that is why we don't need calc
    box: {display: "flex", flexDirection: "column", gap: 2, flexWrap: "nowrap", height: "calc(100vh - 8px)", minWidth: "1010px"},
    body: {display: "flex", alignItems: "center", justifyContent: "center", flexGrow: 1},
}

export function App() {
    const info = useQuery({queryKey: ["info"], queryFn: GeneralApi.info, refetchOnWindowFocus: "always"})
    const show = !info.isPending && !!info.data && info.data.secret.key && info.data.auth.authorised

    return (
        <Box sx={SX.box}>
            <Box>
                <Header company={info.data?.company ?? "VEEGRES"} show={show} auth={info.data?.auth.type ?? AuthType.NONE}/>
            </Box>
            <Box sx={SX.body}>
                <Body info={info}/>
            </Box>
            <Box>
                <Footer tag={info.data?.version.tag ?? "none"} commit={info.data?.version.commit ?? "none"}/>
            </Box>
        </Box>
    );
}
