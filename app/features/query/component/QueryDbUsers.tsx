import { Cancel, CheckCircle, Delete, Edit } from "@mui/icons-material";
import { Box, Button, Chip, TextField } from "@mui/material";
import { useState } from "react";

import { ErrorLeaderMissing } from "../../../shared/component/box/ErrorManual";
import { ErrorSmart } from "../../../shared/component/box/ErrorSmart";
import { NoBox } from "../../../shared/component/box/NoBox";
import { TitleBox } from "../../../shared/component/box/TitleBox";
import { SimpleButton } from "../../../shared/component/button/SimpleButton";
import { FieldRow } from "../../../shared/component/input/FieldRow";
import { SkeletonGroup } from "../../../shared/component/progress/SkeletonGroup";
import { Feature } from "../../Feature";
import { ManageAccess } from "../../management/component/ManageAccess";
import { Role } from "../../node/api/NodeType";
import {
  useRouterDbRoleCreate,
  useRouterDbRoleDelete,
  useRouterDbRoleList,
  useRouterDbRoleUpdate,
} from "../api/QueryHook";
import { Connection, DbRoleItem } from "../api/QueryType";
import { QueryDbAdminPanel, QueryDbAdminSX } from "./QueryDbAdminPanel";

type Props = {
  connection: Connection;
  role: Role;
};

export function QueryDbUsers(props: Props) {
  const { connection, role } = props;
  const list = useRouterDbRoleList(connection);
  const [newName, setNewName] = useState("");
  const [newPassword, setNewPassword] = useState("");

  const resetNew = () => {
    setNewName("");
    setNewPassword("");
  };
  const create = useRouterDbRoleCreate(connection, resetNew);

  if (role !== "leader") return <ErrorLeaderMissing />;

  return (
    <QueryDbAdminPanel
      connection={connection}
      viewFeature={Feature.ViewQueryDbRoleList}
      createFeature={Feature.ManageQueryDbRoleCreate}
      listKind={"role"}
      createLabel={"New role"}
      alert={
        "Create login roles with a password. Attributes such as SUPERUSER are read-only here — use the SQL console for advanced options."
      }
      renderCreate={renderCreateForm()}
    >
      {renderList()}
    </QueryDbAdminPanel>
  );

  function renderCreateForm() {
    return (
      <>
        <FieldRow>
          <TextField
            size={"small"}
            label={"Name"}
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
          />
          <TextField
            size={"small"}
            label={"Password"}
            type={"password"}
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
          />
        </FieldRow>
        <Button
          sx={{ marginTop: 1 }}
          size={"small"}
          variant={"contained"}
          color={"secondary"}
          disabled={!newName.trim() || !newPassword || create.isPending}
          onClick={() =>
            create.mutate({
              connection,
              name: newName.trim(),
              password: newPassword,
            })
          }
        >
          Create role
        </Button>
      </>
    );
  }

  function renderList() {
    if (list.isPending)
      return (
        <SkeletonGroup count={4} width={"100%"} height={52} grow={false} />
      );
    if (list.isError) return <ErrorSmart error={list.error} />;
    if (!list.data || list.data.length === 0)
      return <NoBox text={"There are no roles yet"} />;

    return list.data.map((item) => (
      <QueryDbUserRow key={item.name} connection={connection} role={item} />
    ));
  }
}

type RowProps = {
  connection: Connection;
  role: DbRoleItem;
};

function QueryDbUserRow(props: RowProps) {
  const { connection, role } = props;
  const [edit, setEdit] = useState(false);
  const [password, setPassword] = useState("");
  const protectedRole = isProtectedRole(role.name);

  const update = useRouterDbRoleUpdate(connection, () => {
    setEdit(false);
    setPassword("");
  });
  const remove = useRouterDbRoleDelete(connection);

  return (
    <TitleBox label={role.name} island={true} dense={true} collapsible={false}>
      <Box sx={QueryDbAdminSX.cardHead}>
        <Box sx={QueryDbAdminSX.chips}>
          {protectedRole && (
            <Chip size={"small"} label={"protected"} variant={"outlined"} />
          )}
          {role.superuser && (
            <Chip size={"small"} label={"SUPERUSER"} color={"warning"} />
          )}
          {role.createDb && <Chip size={"small"} label={"CREATEDB"} />}
          {role.login ? (
            <Chip size={"small"} label={"LOGIN"} color={"success"} />
          ) : (
            <Chip size={"small"} label={"NO LOGIN"} variant={"outlined"} />
          )}
        </Box>
        {!protectedRole && !edit && renderActions()}
      </Box>
      {edit && (
        <FieldRow>
          <TextField
            size={"small"}
            label={"New password"}
            type={"password"}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <SimpleButton
            tooltip={"Cancel"}
            onClick={() => setEdit(false)}
            disabled={update.isPending}
          >
            <Cancel fontSize={"small"} />
          </SimpleButton>
          <SimpleButton
            tooltip={"Save password"}
            loading={update.isPending}
            disabled={!password}
            onClick={() =>
              update.mutate({ connection, name: role.name, password })
            }
          >
            <CheckCircle fontSize={"small"} />
          </SimpleButton>
        </FieldRow>
      )}
    </TitleBox>
  );

  function renderActions() {
    return (
      <Box sx={QueryDbAdminSX.cardActions}>
        <ManageAccess feature={Feature.ManageQueryDbRoleUpdate}>
          <SimpleButton
            tooltip={"Change password"}
            onClick={() => setEdit(true)}
            disabled={remove.isPending}
          >
            <Edit fontSize={"small"} />
          </SimpleButton>
        </ManageAccess>
        <ManageAccess feature={Feature.ManageQueryDbRoleDelete}>
          <SimpleButton
            tooltip={"Delete role"}
            loading={remove.isPending}
            onClick={() => remove.mutate({ connection, name: role.name })}
          >
            <Delete fontSize={"small"} />
          </SimpleButton>
        </ManageAccess>
      </Box>
    );
  }
}

function isProtectedRole(name: string) {
  return name === "postgres" || name.startsWith("pg_");
}
