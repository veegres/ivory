import {TitleBox} from "../../../shared/component/box/TitleBox"
import {Options} from "../api/ClusterType"
import {ClusterOptions} from "./ClusterOptions"

type Props = {
    options: Options,
    onUpdate: (options: Options) => void,
}

// ClusterOptionsBox is the cluster options widget as a dialog shows it: folded
// away under its own heading, with the plugin selectors disabled, since inside
// a dialog the plugins come from the cluster list filter rather than the form.
export function ClusterOptionsBox(props: Props) {
    const {options, onUpdate} = props

    return (
        <TitleBox label={"All Configuration"} dense={true}>
            <ClusterOptions options={options} onUpdate={onUpdate} disablePlugins={true}/>
        </TitleBox>
    )
}
