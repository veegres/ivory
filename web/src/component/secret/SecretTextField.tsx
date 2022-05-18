import {TextField} from "@mui/material";
import {InputProps as StandardInputProps} from "@mui/material/Input/Input";

type Props = {
    label: string
    onChange: StandardInputProps['onChange']
}

export function SecretTextField(props: Props) {
    const { label, onChange } = props
    return (
        <TextField
            required
            fullWidth
            margin={"normal"}
            autoComplete={"off"}
            size={"small"}
            label={label}
            onChange={onChange}
       />
    )
}
