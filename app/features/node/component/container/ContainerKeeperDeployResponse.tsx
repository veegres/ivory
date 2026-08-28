import {DialogScreen} from "../../../../shared/component/box/DialogScreen"
import {Logs} from "../../../../shared/component/box/Logs"
import {TitledBox} from "../../../../shared/component/box/TitledBox"

// NOTE: the virtualizer inside Logs needs a pixel height, so this is where it
// and the dialog have to agree: --size-dialog (600) less the screen's padding
// (10 top, 10 bottom), less the box's own frame (its padding, heading and
// gap), less the logs' footer and gap. Under rather than exactly on it - a few
// px of slack costs nothing, overshooting brings back the outer scrollbar this
// screen exists without.
const LOGS_HEIGHT = 500

type Props = {
    logs: string[],
}

// ContainerKeeperDeployResponse is what the deploy left behind - the output of
// the command and of the post script that followed it.
export function ContainerKeeperDeployResponse(props: Props) {
    const {logs} = props

    return (
        <DialogScreen fit={true}>
            <TitledBox title={"Logs"} island={true}>
                <Logs logs={logs} height={LOGS_HEIGHT} auto={false}/>
            </TitledBox>
        </DialogScreen>
    )
}
