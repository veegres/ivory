import {Box} from "@mui/material"
import {ReactNode, useState} from "react"

import {useRouterInfo} from "../../../features/management/api/ManagementHook"
import select from "../../../shared/style/select.module.css"
import {SxPropsMap} from "../../helper/HelperType"
import {randomUnicodeAnimal} from "../../helper/HelperUtils"
import {LinearProgressStateful} from "../progress/LinearProgressStateful"

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", alignItems: "center", maxWidth: "450px",
        width: {xs: "92%", md: "30%"}, minWidth: {xs: 0, md: "450px"}, gap: 2,
    },
    header: {fontSize: "25px", fontWeight: 900, fontFamily: "monospace"},
    animal: {fontSize: "45px", cursor: "pointer"},
    buttons: {display: "flex", gap: 2, width: "100%", justifyContent: "center", alignItems: "center"},
}

type Props = {
    children: ReactNode,
    header: string,
    renderFooter: ReactNode,
    position?: "center" | "start" | "end",
    padding?: string,
}

export function PageStartupBox(props: Props) {
    const {children, header, renderFooter, position = "center", padding = "50px 0px"} = props
    const [animal, setAnimal] = useState(randomUnicodeAnimal())

    const info = useRouterInfo()

    return (
        <Box sx={[SX.box, {justifyContent: position, padding}]}>
            <Box sx={SX.animal} className={select.none} onClick={() => setAnimal(randomUnicodeAnimal())}>{animal}</Box>
            <Box sx={SX.header} className={select.none}>{header}</Box>
            {children}
            <LinearProgressStateful loading={info.isFetching} line color={"inherit"}/>
            <Box sx={SX.buttons}>
                {renderFooter}
            </Box>
        </Box>
    )
}
