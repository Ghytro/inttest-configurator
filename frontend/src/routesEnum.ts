const routesEnum = {
    auth: "/auth",
    projects: "/projects",
    projectPageMask: "/projects/:id",
    users: "/users",
};

function projPageUrl(id: number): string {
    return "/projects/"+id.toString();
}

export { routesEnum, projPageUrl };
