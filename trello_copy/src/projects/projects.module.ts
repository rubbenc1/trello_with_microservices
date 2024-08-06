import { Module } from '@nestjs/common';
import { ProjectsController } from './projects.controller';
import { ProjectsService } from './projects.service';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Project } from './projects.model';
import { UserProjects } from './user-projects.model';
import { User } from 'src/users/users.model';
import { AuthModule } from 'src/auth/auth.module';
import { TasksColumn } from 'src/tasks-column/tasks-column.model';
import { Lists } from 'src/lists/lists.model';
import { Tasks } from 'src/tasks/tasks.model';
import { TasksValuesModule } from 'src/tasks-values/tasks-values.module';


@Module({
  controllers: [ProjectsController],
  providers: [ProjectsService],
  imports: [
    TypeOrmModule.forFeature([Project, User, UserProjects, TasksColumn, Lists, Tasks]),
    AuthModule,
    TasksValuesModule
  ]

})
export class ProjectsModule {}
