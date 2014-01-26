# This cookbook configuration is used to define the tasks required to 
# maintain the diSimplexRepository paper.
#

# Ensure we have an up to date .gitignore
#
gitIgnoreParts = FileList.new( '.gitignore.d/*-gitignore')
file '.gitignore' => gitIgnoreParts do
  puts "cat `ls .gitignore.d/*-gitignore'` > .gitignore"
  system "cat `ls .gitignore.d/*-gitignore` > .gitignore"
end
file '.gitignore' => 'create:default';
task :create => '.gitignore';

add_cookbook("/home/stg/ExpositionGit/tools/latexCookbook")


