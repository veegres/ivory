import {Box} from "@mui/material";
import {ReactNode, useState} from "react";
import {randomUnicodeAnimal} from "../../../app/utils";
import {LinearProgressStateful} from "../progress/LinearProgressStateful";
import select from "../../../style/select.module.css";
import {useRouterInfo} from "../../../api/management/hook";
import {SxPropsMap} from "../../../app/type";

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center",
        width: "30%", minWidth: "500px", gap: 2
    },
    header: {fontSize: '35px', fontWeight: 900, fontFamily: 'monospace', margin: "20px 0", cursor: "pointer"},
    buttons: {display: "flex", gap: 2, width: "100%", justifyContent: "center", alignItems: "center"},
}

type Props = {
    children: ReactNode,
    header: string,
    renderFooter: ReactNode,
}

export function PageStartupBox(props: Props) {
    const {children, header, renderFooter} = props
    const [animal, setAnimal] = useState(randomUnicodeAnimal())

    const info = useRouterInfo()

    return (
        <Box sx={SX.box}>
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
