import {useState} from "react";
import {useForm} from "@mantine/form"
import {Button, Group, Modal, Textarea, TextInput} from "@mantine/core";
import {ENDPOINT,Todo} from "../App";
import {KeyedMutator} from "swr";

function AddTodo({mutate}:{mutate:KeyedMutator<Todo[]>}){
    const [open,setOpen]=useState(false)

    const form=useForm({
        initialValues:{
            name:"",
            description:"",
            audit:{
                created_at:Date.now(),
                updated_at:Date.now(),
            }
        },
    })
   async function createTodo(values:{name:string,description:string}){
        const updated =await fetch(`${ENDPOINT}/todo/create`,{
            method:"POST",
            headers:{
                "Content-Type":"application/json",
                "Email":"random@random.com"
            },
            body:JSON.stringify(values),
        }).then((r)=>r.json());
        mutate(updated)
        form.reset()
       setOpen(false)
    }
    return (
        <>
    <Modal
        opened={open} onClose={()=>setOpen(false)} title="Create Todo">
        <form onSubmit={form.onSubmit(createTodo)}>
            <TextInput
            required
            mb={12}
            label="Todo"
            placeholder="What do you want to do?"
            {...form.getInputProps("name")}
            />
            <Textarea
                required
                mb={12}
                label="Description"
                placeholder="Tell me more.."
                {...form.getInputProps("description")}
            />
            <Button type="submit">Create Todo</Button>
        </form>
    </Modal>
            <Group position="center">
                <Button fullWidth mb={12} onClick={()=>setOpen(true)}>
                    ADD TODO
                </Button>
            </Group>
    </>
    );
}
export default AddTodo