import {Box} from "@mui/material"
import {ReactNode, useState} from "react"

import {useRouterInfo} from "../../../features/management/api/hook"
import select from "../../../shared/style/select.module.css"
import {SxPropsMap} from "../../helper/type"
import {randomUnicodeAnimal} from "../../helper/utils"
import {LinearProgressStateful} from "../progress/LinearProgressStateful"

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", alignItems: "center",
        width: "30%", minWidth: "500px", gap: 2
    },
    header: {fontSize: "35px", fontWeight: 900, fontFamily: "monospace", margin: "20px 0", cursor: "pointer"},
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
    const {children, header, renderFooter, position = "center", padding = "150px 0px"} = props
    const [animal, setAnimal] = useState(randomUnicodeAnimal())

    const info = useRouterInfo()

    return (
        <Box sx={SX.box} justifyContent={position} padding={padding}>
            <Box sx={SX.header} className={select.none} onClick={() => setAnimal(randomUnicodeAnimal())}>
                {header} {animal}
            </Box>
            {children}
            <LinearProgressStateful loading={info.isFetching} line color={"inherit"}/>
            <Box sx={SX.buttons}>
                {renderFooter}
            </Box>
        </Box>
    )
}
