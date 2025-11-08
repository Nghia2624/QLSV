CREATE EXTENSION IF NOT EXISTS "pgcrypto";
--tạo mã code tự động 
CREATE SEQUENCE IF NOT EXISTS student_code_seq START 1;
CREATE SEQUENCE IF NOT EXISTS teacher_code_seq START 1;
CREATE SEQUENCE IF NOT EXISTS class_code_seq START 1;
CREATE SEQUENCE IF NOT EXISTS course_code_seq START 1;

--bảng users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'teacher', 'student')),
    ref_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

--bảng student
CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_code VARCHAR(8) UNIQUE NOT NULL DEFAULT ('SV' || LPAD(nextval('student_code_seq')::text, 3, '0')),
    name VARCHAR(100) NOT NULL,                   
    email VARCHAR(100) UNIQUE NOT NULL,          
    dob VARCHAR(10),                              
    gender VARCHAR(10) CHECK (gender IN ('Nam', 'Nữ')),
    address VARCHAR(255),                     
    phone VARCHAR(20),                        
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS teachers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    teacher_code VARCHAR(8) UNIQUE NOT NULL DEFAULT ('GV' || LPAD(nextval('teacher_code_seq')::text, 3, '0')),
    name VARCHAR(100) NOT NULL,                  
    email VARCHAR(100) UNIQUE NOT NULL,           
    department VARCHAR(100),                      
    address VARCHAR(255),                      
    phone VARCHAR(20),                          
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_code VARCHAR(8) UNIQUE NOT NULL DEFAULT ('MH' || LPAD(nextval('course_code_seq')::text, 3, '0')),
    name VARCHAR(100) NOT NULL,                   
    description TEXT,                             
    credits INTEGER NOT NULL CHECK (credits > 0), 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS classes (

    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_code VARCHAR(8) UNIQUE NOT NULL DEFAULT ('CL' || LPAD(nextval('class_code_seq')::text, 3, '0')),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,  
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,                    
    schedule VARCHAR(100),                         
    semester VARCHAR(20),                   
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS registrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE, 
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, 
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'dropped', 'completed')), 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, class_id)
);

-- Indexes cho bảng users
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

CREATE INDEX IF NOT EXISTS idx_students_email ON students(email);        -- Tìm kiếm theo email
CREATE INDEX IF NOT EXISTS idx_students_created_at ON students(created_at); -- Sắp xếp theo thời gian tạo

-- Indexes cho bảng teachers  
CREATE INDEX IF NOT EXISTS idx_teachers_email ON teachers(email);        -- Tìm kiếm theo email
CREATE INDEX IF NOT EXISTS idx_teachers_department ON teachers(department); -- Tìm kiếm theo khoa

-- Indexes cho bảng courses
CREATE INDEX IF NOT EXISTS idx_courses_name ON courses(name);            -- Tìm kiếm theo tên môn học

-- Indexes cho bảng classes
CREATE INDEX IF NOT EXISTS idx_classes_course_id ON classes(course_id);  -- Tìm kiếm lớp theo môn học
CREATE INDEX IF NOT EXISTS idx_classes_teacher_id ON classes(teacher_id); -- Tìm kiếm lớp theo giảng viên
CREATE INDEX IF NOT EXISTS idx_classes_semester ON classes(semester);    -- Tìm kiếm lớp theo học kỳ

-- Indexes cho bảng registrations
CREATE INDEX IF NOT EXISTS idx_registrations_student_id ON registrations(student_id); -- Tìm kiếm đăng ký theo sinh viên
CREATE INDEX IF NOT EXISTS idx_registrations_class_id ON registrations(class_id);     -- Tìm kiếm đăng ký theo lớp
CREATE INDEX IF NOT EXISTS idx_registrations_status ON registrations(status);         -- Tìm kiếm theo trạng thái

-- Function để tự động cập nhật updated_at khi có thay đổi dữ liệu
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    -- Cập nhật updated_at thành thời gian hiện tại
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Áp dụng triggers cho tất cả bảng
-- Trigger sẽ tự động chạy function update_updated_at_column() trước mỗi lần UPDATE
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_students_updated_at BEFORE UPDATE ON students
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_teachers_updated_at BEFORE UPDATE ON teachers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_courses_updated_at BEFORE UPDATE ON courses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_classes_updated_at BEFORE UPDATE ON classes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_registrations_updated_at BEFORE UPDATE ON registrations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
