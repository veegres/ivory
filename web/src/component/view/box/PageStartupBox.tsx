import {Box} from "@mui/material";
import {ReactNode, useState} from "react";
import {useQuery} from "@tanstack/react-query";
import {generalApi} from "../../../app/api";
import {randomUnicodeAnimal} from "../../../app/utils";
import {LinearProgressStateful} from "../progress/LinearProgressStateful";
import select from "../../../style/select.module.css";
import {SxPropsMap} from "../../../type/common";

const SX: SxPropsMap = {
    box: {
        display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center",
        height: "100%", width: "30%", minWidth: "500px", gap: 2
    },
    header: {fontSize: '35px', fontWeight: 900, fontFamily: 'monospace', margin: "20px 0", cursor: "pointer"},
    buttons: {display: "flex", gap: 4},
}

type Props = {
    children: ReactNode,
    header: string,
    renderFooter: ReactNode,
}

export function PageStartupBox(props: Props) {
    const {children, header, renderFooter} = props
    const [animal, setAnimal] = useState(randomUnicodeAnimal())

    const info = useQuery(["info"], generalApi.info);

    return (
        <Box sx={SX.box}>
            <Box sx={SX.header} className={select.none} onClick={() => setAnimal(randomUnicodeAnimal())}>
                {header} {animal}
            </Box>
            {children}
            <LinearProgressStateful isFetching={info.isFetching} line color={"inherit"}/>
            <Box sx={SX.buttons}>
                {renderFooter}
            </Box>
        </Box>
    )
}
