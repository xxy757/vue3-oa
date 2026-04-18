package utils

import "golang.org/x/crypto/bcrypt"

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// 步骤2：将哈希字节转换回字符串用于存储

// CheckPassword 将明文密码与之前生成的 bcrypt 哈希进行比较，判断是否匹配。
//

	return err == nil
