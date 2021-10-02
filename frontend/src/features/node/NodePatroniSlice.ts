import {createAsyncThunk, createSlice} from '@reduxjs/toolkit';
import {RootState} from '../../app/store';
import {restClient} from "../../app/hooks";

export interface NodePatroni {
    database_system_identifier: string,
    postmaster_start_time: string,
    timeline: number,
    cluster_unlocked: boolean,
    patroni: { scope: string, version: string },
    state: string,
    role: string,
    xlog: {
        received_location: number,
        replayed_timestamp: string,
        paused: boolean,
        replayed_location: number
    }
    server_version: string
}

export interface ClusterListState {
    data?: NodePatroni;
    error?: string
    status: 'idle' | 'loading' | 'failed';
}

const initialState: ClusterListState = {
    status: 'idle',
};

export const incrementAsync = createAsyncThunk(
    'node/fetchPatroni',
    async (node: string, { rejectWithValue }) => {
        try {
            return await restClient<NodePatroni>(`/node/${node}/patroni`)
        } catch (error) {
            return rejectWithValue(error.response.data)
        }
    }
);

export const nodePatroniSlice = createSlice({
    name: 'node/patroni',
    initialState,
    reducers: {},
    extraReducers: (builder) => {
        builder
            .addCase(incrementAsync.pending, (state) => {
                state.status = 'loading';
            })
            .addCase(incrementAsync.fulfilled, (state, action) => {
                state.status = 'idle';
                state.data = action.payload.response;
            })
            .addCase(incrementAsync.rejected, (state, action) => {
                state.status = 'idle';
                state.error = action.error.message;
            });
    },
});

export const selectNodePatroniData = (state: RootState) => state.nodePatroni.data;
