import {Box, List, ThemeIcon} from "@mantine/core";
import './App.css'
import useSWR from "swr";
import AddTodo from "./components/AddTodo"
import {CheckCircleFillIcon} from "@primer/octicons-react";
import {jsx} from "@emotion/react";

export interface Audit{
    created_at:bigint,
    created_by:string,
    updated_at:bigint,
    updated_by:string
}

export interface Todo {
    id:string,
    name:string,
    description:string,
    status:string,
    audit:Audit
}

export const ENDPOINT="http://localhost:8080"
const fetcher = (url:string)=>
    fetch(`${ENDPOINT}/${url}`).then((r)=>r.json());
function App() {
    const {data, mutate}=useSWR<Todo[]>('todo/getall',fetcher);

    async function updateStatus(todo:Todo){
        todo.status = (todo.status === "pending") ? "done" : "pending";
        const updated=await fetch(`${ENDPOINT}/todo/${todo.id}`,{
            headers: {
                'Content-Type': 'application/json'
            },
            method:"PUT",
            body:JSON.stringify(todo)

        }).then((r)=>r.json());
        mutate(updated)
    }
  return (
      <Box sx={(theme)=>({
          padding:"2rem",
          width:"100%",
          maxWidth:"40rem",
          margin:"0 auto",
      })} >
          <List spacing="xs" size="sm" mb={12}>
              {data?.map((todo)=>{
                  return (
                      <List.Item
                          onClick={()=>updateStatus(todo)}
                          key={`todo_list__${todo.id}`}
                      icon={
                          (todo.status !== "pending" && todo.status !== "") ? (<ThemeIcon color="teal" size={24} radius="xl">
                                  <CheckCircleFillIcon size={20}/>
                                  </ThemeIcon>
                          ) : (
                              <ThemeIcon color="grey" size={24} radius="xl">
                                  <CheckCircleFillIcon size={20}/>
                              </ThemeIcon>
                          )
                      }
                      >{todo.name}</List.Item>
                  );
              })}
          </List>
        <AddTodo mutate={mutate}/>
      </Box>
  )
}

export default App
