import type { Meta, StoryObj } from "@storybook/react";
import { StorylineForm } from "./StorylineForm";

const meta = {
  component: StorylineForm,
  tags: ["autodocs"],
} satisfies Meta<typeof StorylineForm>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Create: Story = {
  args: {
    type: "create",
  },
};

export const Edit: Story = {
  args: {
    type: "create",
  },
};
