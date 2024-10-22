
-- +migrate Up
CREATE index average_colors_file_name_index ON average_colors (file_name);

-- +migrate Down
DROP INDEX average_colors_file_name_index;
