import {json} from "@codemirror/lang-json"
import {Box, Skeleton} from "@mui/material"
import ReactCodeMirror from "@uiw/react-codemirror"
import {useEffect, useState} from "react"

import {ErrorKeeperMissing} from "../../../../shared/component/box/ErrorManual"
import {ErrorSmart} from "../../../../shared/component/box/ErrorSmart"
import {TitledBox} from "../../../../shared/component/box/TitledBox"
import {CancelIconButton, CopyIconButton, EditIconButton, SaveIconButton} from "../../../../shared/component/button/IconButtons"
import {SxPropsMap} from "../../../../shared/helper/HelperType"
import {CodeThemes} from "../../../../shared/helper/HelperUtils"
import {useSettings} from "../../../../shared/provider/AppProvider"
import {useSnackbar} from "../../../../shared/provider/SnackbarProvider"
import {Feature} from "../../../Feature"
import {ManageAccess} from "../../../management/component/ManageAccess"
import {useRouterNodeConfig, useRouterNodeConfigUpdate} from "../../api/NodeHook"
import {KeeperOneRequest} from "../../api/NodeType"

const SX: SxPropsMap = {
    input: {flexGrow: 1, borderWidth: "1px", borderStyle: "solid", overflowX: "auto", ">div": {height: "100%"}},
    buttons: {display: "flex"},
}

type Props = {
    req?: KeeperOneRequest,
}

export function KeeperConfig(props: Props) {
    const {req} = props
    const settings = useSettings()
    const snackbar = useSnackbar()
    const [isEditable, setIsEditable] = useState(false)
    const [configState, setConfigState] = useState("")

    const config = useRouterNodeConfig(req)
    const updateConfig = useRouterNodeConfigUpdate(req, () => setIsEditable(false))

    const {data, isPending, isError, error} = config

    useEffect(() => setConfigState(stringify(data)), [data])

    if (!req) return <ErrorKeeperMissing/>
    if (isError) return <ErrorSmart error={error}/>
    if (isPending) return <Skeleton variant={"rectangular"} height={300}/>

    return (
        <TitledBox title={"Config"} island={true} renderActions={renderActions()}>
            <Box sx={[SX.input, {borderColor: isEditable ? "divider" : "transparent"}]}>
                <ReactCodeMirror
                    height={"100%"}
                    width={"100%"}
                    value={configState}
                    editable={isEditable}
                    autoFocus={isEditable}
                    basicSetup={{highlightActiveLine: false, highlightActiveLineGutter: isEditable, highlightSelectionMatches: false}}
                    theme={CodeThemes[settings.theme]}
                    extensions={[json()]}
                    onChange={(value) => setConfigState(value)}
                />
            </Box>
        </TitledBox>
    )

    function renderActions() {
        return (
            <Box sx={SX.buttons}>
                <ManageAccess feature={Feature.ManageNodeKeeperConfigUpdate}>
                    {renderUpdateButtons()}
                </ManageAccess>
                <CopyIconButton placement={"left"} size={30} onClick={handleCopyAll}/>
            </Box>
        )
    }

    function renderUpdateButtons() {
        if (!isEditable) return <EditIconButton placement={"left"} size={30} onClick={() => setIsEditable(true)}/>

        return (
            <>
                <SaveIconButton placement={"left"} size={30} loading={updateConfig.isPending} onClick={handleUpdate}/>
                <CancelIconButton placement={"left"} size={30} disabled={updateConfig.isPending} onClick={handleCancel}/>
            </>
        )
    }

    function handleCopyAll() {
        const currentConfig = configState ? configState : stringify(data)
        navigator.clipboard.writeText(currentConfig).then(() => {
            snackbar("Config copied to clipboard!", "info")
        })
    }

    function handleCancel() {
        setIsEditable(false)
        setConfigState(stringify(data))
    }

    function handleUpdate() {
        if (configState && req) {
            updateConfig.mutate({...req, body: JSON.parse(configState)})
        }
    }

    function stringify(json?: any) {
        return JSON.stringify(json, null, 2)
    }
}
