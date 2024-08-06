import { Module } from '@nestjs/common';
import { ListsService } from './lists.service';
import { ListsController } from './lists.controller';
import { Project } from 'src/projects/projects.model';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Lists } from './lists.model';
import { UserProjects } from 'src/projects/user-projects.model';
import { AuthModule } from 'src/auth/auth.module';
import { Tasks } from 'src/tasks/tasks.model';
import { TasksValuesModule } from 'src/tasks-values/tasks-values.module';

@Module({
  providers: [ListsService],
  controllers: [ListsController],
  imports: [
    TypeOrmModule.forFeature([Project, Lists, UserProjects, Tasks]),
    AuthModule,
    TasksValuesModule
  ]
})
export class ListsModule {}
