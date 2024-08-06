import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { CreateProject } from './dto/create-project.dto';
import { Repository, In } from 'typeorm';
import { Project } from './projects.model';
import { UserProjects } from './user-projects.model';
import { User } from 'src/users/users.model';
import { Tasks } from 'src/tasks/tasks.model';
import { Lists } from 'src/lists/lists.model';
import { TasksValuesService } from 'src/tasks-values/tasks-values.service';

@Injectable()
export class ProjectsService {
    constructor(
        @InjectRepository(Project) private projectRepository: Repository<Project>,
        @InjectRepository(UserProjects) private userProjectsRepository: Repository <UserProjects>,
        @InjectRepository(User) private userRepository: Repository<User>,
        @InjectRepository(Lists) private listsRepository: Repository<Lists>,
        @InjectRepository(Tasks) private tasksRepository: Repository<Tasks>,
        private tasksValuesService: TasksValuesService
    ){}


    async createProject(name: CreateProject, userId: number){
        const user = await this.userRepository.findOneBy({id: userId})
        if(!user){
            throw new NotFoundException(`User with ID ${userId} not found`);
        }
        const project = this.projectRepository.create(name)
        await this.projectRepository.save(project)

        const userProject = this.userProjectsRepository.create({userId, projectId: project.id});
        await this.userProjectsRepository.save(userProject);
        return project;
    }

    async deleteProject(id: number) {
        // Find the project to delete
        const projectToDelete = await this.projectRepository.findOne({
            where: { id },
            relations: ['lists', 'lists.tasks']
        });

        if (!projectToDelete) {
            throw new NotFoundException(`Project with ID ${id} not found`);
        }

        const { lists } = projectToDelete;

        // Delete tasks and their values for each list
        for (const list of lists) {
            const tasks = await this.tasksRepository.find({ where: { listId: list.id } });

            for (const task of tasks) {
                await this.tasksValuesService.deleteTaskOfAllValues(task.id);
            }

            await this.tasksRepository.delete({ listId: list.id });
        }

        // Delete all lists in the project
        await this.listsRepository.delete({ projectId: id });

        // Delete the project
        await this.projectRepository.delete(id);

        // Optionally, remove user-project associations
        await this.userProjectsRepository.delete({ projectId: id });
    }

    async updateProject(id: number, name: CreateProject){
        const result = await this.projectRepository.update(id, name);
        if(result.affected === 0){
            throw new NotFoundException(`Project with ID ${id} not found`);
        }
        const project = await this.projectRepository.findOneBy({id});
        return project;
    }
    async getProject(id:number){
        const result = await this.projectRepository.findOne({
            where:{id},
            relations:['lists']
        });
        if(!result){
            throw new NotFoundException(`Project with ID ${id} not found`);
        }
        return result
    }
    async getUserProjects(userId: number){
        const userProjects = await this.userProjectsRepository.find({where: {userId}});
        const projectIds = userProjects.map((item)=> item.projectId);
        if (projectIds.length === 0) {
            return [];
        }
        const projectsName = this.projectRepository.find({
            where:{id: In (projectIds)},
            relations: ['lists'],
        })
        return projectsName;
    }
    async getAllProjects(){
        const projects = await this.projectRepository.find({ relations: ['lists', 'lists.tasks'] });
        return projects
    }
}
