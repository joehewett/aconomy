import React from "react"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog"
import { AgentTurn } from "../Game"


interface PromptModalProps {
  agentTurn: AgentTurn;
}


export const PromptModal = ({ agentTurn }: PromptModalProps) =>
  <Dialog>
    <DialogTrigger>Prompt</DialogTrigger>
    <DialogContent className="max-w-[700px] max-h-[900px] overflow-auto">
      <DialogHeader>
        <DialogTitle>Prompt for Agent {agentTurn.AgentID} on turn {agentTurn.Turn}</DialogTitle>
        <DialogDescription>
          <div className="space-y-3">

            {agentTurn.FullPrompt.map((prompt, idx) => (
              <div key={idx}>
                <p className="text-lg font-semibold">{prompt.role}</p>
                <pre>
                  <p>{prompt.content}</p>
                </pre>
              </div>
            ))}
          </div>
        </DialogDescription>
      </DialogHeader>
    </DialogContent>
  </Dialog>
