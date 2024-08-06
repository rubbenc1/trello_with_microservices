import { ApiProperty } from "@nestjs/swagger";
import { TasksColumn } from "src/tasks-column/tasks-column.model";
import { Column, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn, Unique } from "typeorm";


@Entity({name: 'tasks-column-options'})
@Unique(['columnId', 'optionName'])
export class TasksColumnOptions{

    @ApiProperty({example: '1', description: 'option id'})
    @PrimaryGeneratedColumn()
    optionId: number;

    @ApiProperty({example: 'High', description: 'column name'})
    @Column()
    optionName: string;

    @ManyToOne(()=>TasksColumn, (column)=>column.options, {onDelete: 'CASCADE'})
    @JoinColumn({name:'columnId'})
    column: TasksColumn;
    @Column()
    columnId: number;
    
}