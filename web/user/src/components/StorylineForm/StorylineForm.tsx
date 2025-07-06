export type StorylineFormProps = Readonly<{
  type: "create" | "edit";
}>;

export const StorylineForm: React.FC<StorylineFormProps> = ({ type }) => {
  console.log("StorylineForm, type = ", type);
  return <>form</>;
};
